/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package didanchorinfo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/trustbloc/edge-core/pkg/log"
	"github.com/trustbloc/sidetree-core-go/pkg/api/operation"
	"github.com/trustbloc/sidetree-core-go/pkg/api/protocol"
	"github.com/trustbloc/sidetree-core-go/pkg/docutil"

	"github.com/trustbloc/orb/pkg/didanchor"
	"github.com/trustbloc/orb/pkg/resolver/resource/registry"
)

var logger = log.New("did-anchor-info")

const minDidParts = 4

// ErrDataNotFound describes data not found error.
var ErrDataNotFound = errors.New("data not found")

// DidAnchorInfo retrieves the latest anchor CID and anchor origin for this did (suffix).
type DidAnchorInfo struct {
	namespace    string
	didAnchors   didAnchorProvider
	opsProcessor operationProcessor
}

// didAnchorProvider interface provides access to latest anchor for suffix.
type didAnchorProvider interface {
	Get(suffix string) (string, error)
}

// operationProcessor is an interface which resolves the document based on the suffix.
type operationProcessor interface {
	Resolve(uniqueSuffix string, additionalOps ...*operation.AnchoredOperation) (*protocol.ResolutionModel, error)
}

// New returns a new DidAnchorInfo.
func New(namespace string, didAnchors didAnchorProvider, opsProcessor operationProcessor) *DidAnchorInfo {
	return &DidAnchorInfo{
		namespace:    namespace,
		didAnchors:   didAnchors,
		opsProcessor: opsProcessor,
	}
}

// GetResourceInfo retrieves anchoring info for did(suffix).
func (h *DidAnchorInfo) GetResourceInfo(did string) (registry.Metadata, error) {
	suffix, err := getSuffix(did)
	if err != nil {
		return nil, err
	}

	anchor, err := h.didAnchors.Get(suffix)
	if err != nil {
		if errors.Is(err, didanchor.ErrDataNotFound) {
			return nil, ErrDataNotFound
		}

		return nil, fmt.Errorf("failed to retrieve anchor for suffix[%s]: %w", suffix, err)
	}

	resolutionResult, err := h.opsProcessor.Resolve(suffix)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve operations for suffix[%s]: %w", suffix, err)
	}

	info := make(registry.Metadata)
	info[registry.AnchorURIProperty] = anchor
	info[registry.AnchorOriginProperty] = resolutionResult.AnchorOrigin
	info[registry.CanonicalReferenceProperty] = resolutionResult.CanonicalReference

	logger.Debugf("latest anchor info for suffix[%s]: %+v", suffix, info)

	return info, nil
}

// Accept will accept/reject processing of this did(resource id).
func (h *DidAnchorInfo) Accept(did string) bool {
	parts := strings.Split(did, docutil.NamespaceDelimiter)
	if len(parts) < minDidParts {
		return false
	}

	return strings.HasPrefix(did, h.namespace+docutil.NamespaceDelimiter)
}

func getSuffix(did string) (string, error) {
	lastDelimiter := strings.LastIndex(did, docutil.NamespaceDelimiter)

	adjustedPos := lastDelimiter + 1
	if adjustedPos >= len(did) {
		return "", errors.New("did suffix is empty")
	}

	return did[adjustedPos:], nil
}
