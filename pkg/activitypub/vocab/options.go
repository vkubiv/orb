/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package vocab

import (
	"net/url"
	"time"
)

// Options holds all of the options for building an ActivityPub object.
type Options struct {
	Context      []Context
	ID           *url.URL
	URL          []*url.URL
	To           []*url.URL
	Published    *time.Time
	StartTime    *time.Time
	EndTime      *time.Time
	Types        []Type
	CID          string
	InReplyTo    *url.URL
	Attachment   []*ObjectProperty
	AttributedTo *url.URL
	Generator    string
	Tag          []*TagProperty
	Link         *LinkType

	ObjectPropertyOptions
	CollectionOptions
	ActivityOptions
	ActorOptions
	PublicKeyOptions
	AnchorEventOptions
}

// Opt is an for an object, activity, etc.
type Opt func(opts *Options)

// NewOptions returns an Options struct which is populated with the provided options.
func NewOptions(opts ...Opt) *Options {
	options := &Options{}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

// WithContext sets the 'context' property on the object.
func WithContext(context ...Context) Opt {
	return func(opts *Options) {
		opts.Context = context
	}
}

// WithID sets the 'id' property on the object.
func WithID(id *url.URL) Opt {
	return func(opts *Options) {
		opts.ID = id
	}
}

// WithURL sets the 'url' property on the object.
func WithURL(u ...*url.URL) Opt {
	return func(opts *Options) {
		opts.URL = append(opts.URL, u...)
	}
}

// WithCID sets the 'cid' property on the object.
func WithCID(cid string) Opt {
	return func(opts *Options) {
		opts.CID = cid
	}
}

// WithTo sets the "to" property on the object.
func WithTo(to ...*url.URL) Opt {
	return func(opts *Options) {
		opts.To = append(opts.To, to...)
	}
}

// WithType sets tye 'type' property on the object.
func WithType(t ...Type) Opt {
	return func(opts *Options) {
		opts.Types = t
	}
}

// WithPublishedTime sets the 'publishedTime' property on the object.
func WithPublishedTime(t *time.Time) Opt {
	return func(opts *Options) {
		opts.Published = t
	}
}

// WithStartTime sets the 'startTime' property on the object.
func WithStartTime(t *time.Time) Opt {
	return func(opts *Options) {
		opts.StartTime = t
	}
}

// WithEndTime sets the 'endTime' property on the object.
func WithEndTime(t *time.Time) Opt {
	return func(opts *Options) {
		opts.EndTime = t
	}
}

// WithInReplyTo sets the 'inReplyTo' property on the object.
func WithInReplyTo(id *url.URL) Opt {
	return func(opts *Options) {
		opts.InReplyTo = id
	}
}

// WithAttachment sets the 'attachment' property on the object.
func WithAttachment(obj ...*ObjectProperty) Opt {
	return func(opts *Options) {
		opts.Attachment = append(opts.Attachment, obj...)
	}
}

// WithAttributedTo sets the 'attributedTo' property on the object.
func WithAttributedTo(u *url.URL) Opt {
	return func(opts *Options) {
		opts.AttributedTo = u
	}
}

// WithGenerator sets the 'generator' property on the object.
func WithGenerator(generator string) Opt {
	return func(opts *Options) {
		opts.Generator = generator
	}
}

// WithTag sets the 'tag' property on the object.
func WithTag(tag *TagProperty) Opt {
	return func(opts *Options) {
		if tag != nil {
			opts.Tag = append(opts.Tag, tag)
		}
	}
}

// WithLink sets the 'link' property on the object.
func WithLink(link *LinkType) Opt {
	return func(opts *Options) {
		opts.Link = link
	}
}

// CollectionOptions holds the options for a Collection or OrderedCollection.
type CollectionOptions struct {
	TotalItems int
	First      *url.URL
	Last       *url.URL
	Current    *url.URL
	PartOf     *url.URL
	Next       *url.URL
	Prev       *url.URL
}

// WithTotalItems sets the 'totalItems' property on the collection or ordered collection.
func WithTotalItems(totalItems int) Opt {
	return func(opts *Options) {
		opts.TotalItems = totalItems
	}
}

// WithFirst sets the 'first' property on the collection or ordered collection.
func WithFirst(first *url.URL) Opt {
	return func(opts *Options) {
		opts.First = first
	}
}

// WithLast sets the 'last' property on the collection or ordered collection.
func WithLast(last *url.URL) Opt {
	return func(opts *Options) {
		opts.Last = last
	}
}

// WithCurrent sets the 'current' property on the collection or ordered collection.
func WithCurrent(current *url.URL) Opt {
	return func(opts *Options) {
		opts.Current = current
	}
}

// WithPartOf sets the 'partOf' property on a collection page or ordered collection page.
func WithPartOf(partOf *url.URL) Opt {
	return func(opts *Options) {
		opts.PartOf = partOf
	}
}

// WithNext sets the 'next' property on a collection page or ordered collection page.
func WithNext(next *url.URL) Opt {
	return func(opts *Options) {
		opts.Next = next
	}
}

// WithPrev sets the 'prev' property on a collection page or ordered collection page.
func WithPrev(prev *url.URL) Opt {
	return func(opts *Options) {
		opts.Prev = prev
	}
}

// ObjectPropertyOptions holds options for an 'object' property.
type ObjectPropertyOptions struct {
	Iri               *url.URL
	Object            *ObjectType
	Collection        *CollectionType
	OrderedCollection *OrderedCollectionType
	Activity          *ActivityType
}

// WithIRI sets the 'object' property to an IRI.
func WithIRI(iri *url.URL) Opt {
	return func(opts *Options) {
		opts.Iri = iri
	}
}

// WithObject sets the 'object' property to an embedded object.
func WithObject(obj *ObjectType) Opt {
	return func(opts *Options) {
		opts.Object = obj
	}
}

// WithCollection sets the 'object' property to an embedded collection.
func WithCollection(coll *CollectionType) Opt {
	return func(opts *Options) {
		opts.Collection = coll
	}
}

// WithOrderedCollection sets the 'object' property to an embedded ordered collection.
func WithOrderedCollection(coll *OrderedCollectionType) Opt {
	return func(opts *Options) {
		opts.OrderedCollection = coll
	}
}

// WithActivity sets the 'object' property to an embedded activity.
func WithActivity(activity *ActivityType) Opt {
	return func(opts *Options) {
		opts.Activity = activity
	}
}

// ActivityOptions holds the options for an Activity.
type ActivityOptions struct {
	Result *ObjectProperty
	Actor  *url.URL
	Target *ObjectProperty
}

// WithActor sets the 'actor' property on the activity.
func WithActor(actor *url.URL) Opt {
	return func(opts *Options) {
		opts.Actor = actor
	}
}

// WithTarget sets the 'target' property on the activity.
func WithTarget(target *ObjectProperty) Opt {
	return func(opts *Options) {
		opts.Target = target
	}
}

// WithResult sets the 'result' property on the activity.
func WithResult(result *ObjectProperty) Opt {
	return func(opts *Options) {
		opts.Result = result
	}
}

// ActorOptions holds the options for an Activity.
type ActorOptions struct {
	PublicKey  *PublicKeyType
	Inbox      *url.URL
	Outbox     *url.URL
	Followers  *url.URL
	Following  *url.URL
	Witnesses  *url.URL
	Witnessing *url.URL
	Liked      *url.URL
	Likes      *url.URL
	Shares     *url.URL
}

// WithPublicKey sets the 'publicKey' property on the actor.
func WithPublicKey(publicKey *PublicKeyType) Opt {
	return func(opts *Options) {
		opts.PublicKey = publicKey
	}
}

// WithInbox sets the 'inbox' property on the actor.
func WithInbox(inbox *url.URL) Opt {
	return func(opts *Options) {
		opts.Inbox = inbox
	}
}

// WithOutbox sets the 'outbox' property on the actor.
func WithOutbox(outbox *url.URL) Opt {
	return func(opts *Options) {
		opts.Outbox = outbox
	}
}

// WithFollowers sets the 'followers' property on the actor.
func WithFollowers(followers *url.URL) Opt {
	return func(opts *Options) {
		opts.Followers = followers
	}
}

// WithFollowing sets the 'following' property on the actor.
func WithFollowing(following *url.URL) Opt {
	return func(opts *Options) {
		opts.Following = following
	}
}

// WithWitnesses sets the 'witnesses' property on the actor.
func WithWitnesses(witnesses *url.URL) Opt {
	return func(opts *Options) {
		opts.Witnesses = witnesses
	}
}

// WithWitnessing sets the 'witnessing' property on the actor.
func WithWitnessing(witnessing *url.URL) Opt {
	return func(opts *Options) {
		opts.Witnessing = witnessing
	}
}

// WithLiked sets the 'liked' property on the actor.
func WithLiked(liked *url.URL) Opt {
	return func(opts *Options) {
		opts.Liked = liked
	}
}

// WithLikes sets the 'likes' property on the actor.
func WithLikes(likes *url.URL) Opt {
	return func(opts *Options) {
		opts.Likes = likes
	}
}

// WithShares sets the 'shares' property on the actor.
func WithShares(shares *url.URL) Opt {
	return func(opts *Options) {
		opts.Shares = shares
	}
}

// PublicKeyOptions holds the options for a Public Key.
type PublicKeyOptions struct {
	Owner        *url.URL
	PublicKeyPem string
}

// WithOwner sets the 'owner' property on the public key.
func WithOwner(owner *url.URL) Opt {
	return func(opts *Options) {
		opts.Owner = owner
	}
}

// WithPublicKeyPem sets the 'publicKeyPem' property on the public key.
func WithPublicKeyPem(pem string) Opt {
	return func(opts *Options) {
		opts.PublicKeyPem = pem
	}
}

func getContexts(options *Options, contexts ...Context) []Context {
	return append(options.Context, contexts...)
}

// AnchorEventOptions holds the options for an AnchorEvent.
type AnchorEventOptions struct {
	Index        *url.URL
	Parent       []*url.URL
	AnchorObject *AnchorObjectType
	AnchorEvent  *AnchorEventType
}

// WithIndex sets the 'index' property.
func WithIndex(anchors *url.URL) Opt {
	return func(opts *Options) {
		opts.Index = anchors
	}
}

// WithParent sets the 'parent' property.
func WithParent(parent ...*url.URL) Opt {
	return func(opts *Options) {
		opts.Parent = parent
	}
}

// WithAnchorEvent sets the 'object' property.
func WithAnchorEvent(anchorEvt *AnchorEventType) Opt {
	return func(opts *Options) {
		opts.AnchorEvent = anchorEvt
	}
}

// WithAnchorObject sets the 'object' property to an AnchorObject.
func WithAnchorObject(anchorObj *AnchorObjectType) Opt {
	return func(opts *Options) {
		opts.AnchorObject = anchorObj
	}
}
