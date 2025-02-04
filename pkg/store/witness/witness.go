/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package witness

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/aries-framework-go/spi/storage"
	"github.com/trustbloc/edge-core/pkg/log"

	"github.com/trustbloc/orb/pkg/anchor/witness/proof"
	orberrors "github.com/trustbloc/orb/pkg/errors"
	"github.com/trustbloc/orb/pkg/store/expiry"
)

const (
	namespace = "witness"

	anchorIndex   = "anchorID"
	expiryTagName = "ExpiryTime"
)

var logger = log.New("witness-store")

type updateWitnessProofFnc func(wf *proof.WitnessProof)

// New creates new anchor witness store.
func New(provider storage.Provider, expiryService *expiry.Service, expiryPeriod time.Duration) (*Store, error) {
	store, err := provider.OpenStore(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to open anchor witness store: %w", err)
	}

	s := &Store{
		store:        store,
		expiryPeriod: expiryPeriod,
	}

	err = provider.SetStoreConfig(namespace, storage.StoreConfiguration{TagNames: []string{anchorIndex, expiryTagName}})
	if err != nil {
		return nil, fmt.Errorf("failed to set store configuration: %w", err)
	}

	expiryService.Register(store, expiryTagName, namespace, expiry.WithExpiryHandler(s))

	return s, nil
}

// Store is db implementation of anchor witness store.
type Store struct {
	store        storage.Store
	expiryPeriod time.Duration
}

// Put saves witnesses into anchor witness store.
func (s *Store) Put(anchorID string, witnesses []*proof.Witness) error {
	operations := make([]storage.Operation, len(witnesses))

	anchorIDEncoded := base64.RawURLEncoding.EncodeToString([]byte(anchorID))

	putOptions := &storage.PutOptions{IsNewKey: true}

	for i, w := range witnesses {
		value, err := json.Marshal(w)
		if err != nil {
			return fmt.Errorf("failed to marshal anchor witness: %w", err)
		}

		logger.Debugf("adding %s witness to storage batch: %s", w.Type, w.URI)

		op := storage.Operation{
			Key:   uuid.New().String(),
			Value: value,
			Tags: []storage.Tag{
				{
					Name:  anchorIndex,
					Value: anchorIDEncoded,
				},
				{
					Name:  expiryTagName,
					Value: fmt.Sprintf("%d", time.Now().Add(s.expiryPeriod).Unix()),
				},
			},
			PutOptions: putOptions,
		}

		operations[i] = op
	}

	err := s.store.Batch(operations)
	if err != nil {
		return orberrors.NewTransient(fmt.Errorf("failed to store witnesses for anchorID[%s]: %w", anchorID, err))
	}

	logger.Debugf("stored %d witnesses for anchorID[%s]", len(witnesses), anchorID)

	return nil
}

// Delete deletes all witnesses associated with anchor ID.
func (s *Store) Delete(anchorID string) error {
	var err error

	anchorIDEncoded := base64.RawURLEncoding.EncodeToString([]byte(anchorID))
	query := fmt.Sprintf("%s:%s", anchorIndex, anchorIDEncoded)

	iter, err := s.store.Query(query)
	if err != nil {
		return orberrors.NewTransient(fmt.Errorf("failed to get witnesses for[%s]: %w", query, err))
	}

	defer func() {
		err = iter.Close()
		if err != nil {
			logger.Errorf("failed to close iterator: %s", err.Error())
		}
	}()

	ok, err := iter.Next()
	if err != nil {
		return fmt.Errorf("iterator error for anchorID[%s] : %w", anchorID, err)
	}

	if !ok {
		logger.Debugf("no witnesses to delete for anchorID[%s], nothing to do", anchorID)

		return nil
	}

	var witnessKeys []string

	for ok {
		var key string

		key, err = iter.Key()
		if err != nil {
			return fmt.Errorf("failed to get iterator value for anchorID[%s]: %w", anchorID, err)
		}

		witnessKeys = append(witnessKeys, key)

		ok, err = iter.Next()
		if err != nil {
			return fmt.Errorf("iterator error for anchorID[%s] : %w", anchorID, err)
		}
	}

	operations := make([]storage.Operation, len(witnessKeys))

	for i, k := range witnessKeys {
		operations[i] = storage.Operation{Key: k}
	}

	err = s.store.Batch(operations)
	if err != nil {
		return orberrors.NewTransient(fmt.Errorf("failed to delete witnesses for anchorID[%s]: %w", anchorID, err))
	}

	logger.Debugf("deleted %d witnesses for anchorID[%s]", len(witnessKeys), anchorID)

	return nil
}

// Get retrieves witnesses for the given anchor id.
func (s *Store) Get(anchorID string) ([]*proof.WitnessProof, error) {
	var err error

	anchorIDEncoded := base64.RawURLEncoding.EncodeToString([]byte(anchorID))

	query := fmt.Sprintf("%s:%s", anchorIndex, anchorIDEncoded)

	iter, err := s.store.Query(query)
	if err != nil {
		return nil, orberrors.NewTransient(fmt.Errorf("failed to get witnesses for[%s]: %w", query, err))
	}

	defer func() {
		err = iter.Close()
		if err != nil {
			logger.Errorf("failed to close iterator: %s", err.Error())
		}
	}()

	ok, err := iter.Next()
	if err != nil {
		return nil, orberrors.NewTransient(fmt.Errorf("iterator error for anchorID[%s] : %w", anchorID, err))
	}

	var witnesses []*proof.WitnessProof

	for ok {
		var value []byte

		value, err = iter.Value()
		if err != nil {
			return nil, orberrors.NewTransient(fmt.Errorf("failed to get iterator value for anchorID[%s]: %w",
				anchorID, err))
		}

		var witness proof.WitnessProof

		err = json.Unmarshal(value, &witness)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal anchor witness from store value for anchorID[%s]: %w",
				anchorID, err)
		}

		witnesses = append(witnesses, &witness)

		ok, err = iter.Next()
		if err != nil {
			return nil, orberrors.NewTransient(fmt.Errorf("iterator error for anchorID[%s] : %w", anchorID, err))
		}
	}

	logger.Debugf("retrieved %d witnesses for anchorID[%s]", len(witnesses), anchorID)

	if len(witnesses) == 0 {
		return nil, fmt.Errorf("anchorID[%s] not found in the store", anchorID)
	}

	return witnesses, nil
}

// AddProof adds proof for anchor id and witness.
func (s *Store) AddProof(anchorID string, witness *url.URL, p []byte) error {
	return s.updateWitnessProof(anchorID, []*url.URL{witness}, func(wf *proof.WitnessProof) {
		wf.Proof = p
	})
}

func (s *Store) updateWitnessProof(anchorID string, witnesses []*url.URL, updateFnc updateWitnessProofFnc) error { //nolint:funlen,gocyclo,cyclop,lll
	anchorIDEncoded := base64.RawURLEncoding.EncodeToString([]byte(anchorID))

	query := fmt.Sprintf("%s:%s", anchorIndex, anchorIDEncoded)

	iter, err := s.store.Query(query)
	if err != nil {
		return orberrors.NewTransient(fmt.Errorf("failed to get witnesses for[%s]: %w", query, err))
	}

	defer func() {
		err = iter.Close()
		if err != nil {
			logger.Errorf("failed to close iterator: %s", err.Error())
		}
	}()

	ok, err := iter.Next()
	if err != nil {
		return fmt.Errorf("iterator error for anchorID[%s] : %w", anchorID, err)
	}

	updatedNo := 0

	witnessesMap := getWitnessesMap(witnesses)

	for ok {
		var value []byte

		value, err = iter.Value()
		if err != nil {
			return fmt.Errorf("failed to get iterator value for anchorID[%s]: %w", anchorID, err)
		}

		var w proof.WitnessProof

		err = json.Unmarshal(value, &w)
		if err != nil {
			return fmt.Errorf("failed to unmarshal anchor witness from store value for anchorID[%s]: %w",
				anchorID, err)
		}

		if _, ok = witnessesMap[w.URI.String()]; ok {
			var key string

			key, err = iter.Key()
			if err != nil {
				return fmt.Errorf("failed to get key for anchorID[%s] and witness[%s]: %w",
					anchorID, w.URI.String(), err)
			}

			updateFnc(&w)

			witnessProofBytes, marshalErr := json.Marshal(w)
			if marshalErr != nil {
				return fmt.Errorf("failed to marshal witness[%s] proof for anchorID[%s]: %w", w.URI, anchorID, marshalErr)
			}

			err = s.store.Put(key, witnessProofBytes, storage.Tag{Name: anchorIndex, Value: anchorIDEncoded})
			if err != nil {
				return orberrors.NewTransient(fmt.Errorf("failed to add proof for anchorID[%s] and witness[%s]: %w",
					anchorID, w.URI.String(), err))
			}

			updatedNo++

			logger.Debugf("updated witness proof for anchorID[%s] and witness[%s]", anchorID, w.URI.String())
		}

		ok, err = iter.Next()
		if err != nil {
			return fmt.Errorf("iterator error for anchorID[%s] : %w", anchorID, err)
		}
	}

	if updatedNo == 0 {
		return fmt.Errorf("witness%s not found for anchorID[%s]", witnesses, anchorID)
	}

	return nil
}

// UpdateWitnessSelection updates witness selection flag.
func (s *Store) UpdateWitnessSelection(anchorID string, witnesses []*url.URL, selected bool) error {
	return s.updateWitnessProof(anchorID, witnesses, func(wf *proof.WitnessProof) {
		wf.Selected = selected
	})
}

// HandleExpiredKeys is expired keys inspector/handler.
func (s *Store) HandleExpiredKeys(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	uniqueAnchors := make(map[string]bool)

	for _, key := range keys {
		tags, err := s.store.GetTags(key)
		if err != nil {
			logger.Errorf("get tags for expired key[%s]: %s", key, err)

			return nil
		}

		for _, tag := range tags {
			if tag.Name == anchorIndex {
				anchor, err := base64.RawURLEncoding.DecodeString(tag.Value)
				if err != nil {
					logger.Errorf("failed to decode encoded anchor[%s]: %s", tag.Value, err)

					return nil
				}

				uniqueAnchors[string(anchor)] = true
			}
		}
	}

	anchors := make([]string, 0, len(uniqueAnchors))
	for a := range uniqueAnchors {
		anchors = append(anchors, a)
	}

	logger.Errorf("failed to process anchors: %s", anchors)

	return nil
}

func getWitnessesMap(witnesses []*url.URL) map[string]bool {
	witnessesMap := make(map[string]bool)

	for _, w := range witnesses {
		_, ok := witnessesMap[w.String()]
		if !ok {
			witnessesMap[w.String()] = true
		}
	}

	return witnessesMap
}
