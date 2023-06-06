package controller

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/lefinal/meh"
	"github.com/lefinal/nulls"
	"github.com/mobile-directing-system/mds-server/services/go/logistics-svc/store"
	"github.com/mobile-directing-system/mds-server/services/go/shared/pagination"
	"github.com/mobile-directing-system/mds-server/services/go/shared/pgutil"
	"github.com/mobile-directing-system/mds-server/services/go/shared/search"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// Controller manages all operations regarding logistics.
type Controller struct {
	Logger   *zap.Logger
	DB       pgutil.DBTxSupplier
	Store    Store
	Notifier Notifier
}

// Run the controller for periodic checks, etc.
func (c *Controller) Run(lifetime context.Context) error {
	eg, egCtx := errgroup.WithContext(lifetime)
	eg.Go(func() error {
		return meh.NilOrWrap(c.runPeriodicDeliveryChecks(egCtx), "run periodic delivery checks", nil)
	})
	return eg.Wait()
}

// Store for persistence.
type Store interface {
	// ChannelsByAddressBookEntry retrieves all channels for the address book entry
	// with the given id.
	ChannelsByAddressBookEntry(ctx context.Context, tx pgx.Tx, entryID uuid.UUID) ([]store.Channel, error)
	// AssureAddressBookEntryExists makes sure that the address book entry with the
	// given id exists.
	AssureAddressBookEntryExists(ctx context.Context, tx pgx.Tx, entryID uuid.UUID) error
	// DeleteChannelWithDetailsByID deletes the channel with the given id and type.
	// This is meant to be used as a "shortcut" for clearing channel details as
	// well. This is why we expect the store.ChannelType as well without querying it
	// again.
	DeleteChannelWithDetailsByID(ctx context.Context, tx pgx.Tx, channelID uuid.UUID, channelType store.ChannelType) error
	// UpdateChannelsByEntry clears and recreates channels for the entry with the
	// given id.
	//
	// Warning: No entry existence checks are performed!
	UpdateChannelsByEntry(ctx context.Context, tx pgx.Tx, entryID uuid.UUID, newChannels []store.Channel) error
	// AddressBookEntries retrieves a paginated store.AddressBookEntryDetailed list
	// using the given store.AddressBookEntryFilters and pagination.Params.
	AddressBookEntries(ctx context.Context, tx pgx.Tx, filters store.AddressBookEntryFilters,
		paginationParams pagination.Params) (pagination.Paginated[store.AddressBookEntryDetailed], error)
	// AddressBookEntryByID retrieves the store.AddressBookEntryDetailed with the
	// given id. If visible-by is given, an meh.ErrNotFound will be returned, if the
	// entry is associated with a user, that is not part of any operation, the
	// client (visible-by) is part of.
	AddressBookEntryByID(ctx context.Context, tx pgx.Tx, entryID uuid.UUID, visibleBy uuid.NullUUID) (store.AddressBookEntryDetailed, error)
	// CreateAddressBookEntry creates the given store.AddressBookEntry.
	CreateAddressBookEntry(ctx context.Context, tx pgx.Tx, entry store.AddressBookEntry) (store.AddressBookEntryDetailed, error)
	// UpdateAddressBookEntry updates the given store.AddressBookEntry, identified
	// by its id.
	UpdateAddressBookEntry(ctx context.Context, tx pgx.Tx, entry store.AddressBookEntry) error
	// DeleteAddressBookEntryByID deletes the address book entry with the given id.
	DeleteAddressBookEntryByID(ctx context.Context, tx pgx.Tx, entryID uuid.UUID) error
	// CreateGroup creates the given store.Group.
	CreateGroup(ctx context.Context, tx pgx.Tx, create store.Group) error
	// UpdateGroup updates the given store.Group, identified by its id.
	UpdateGroup(ctx context.Context, tx pgx.Tx, update store.Group) error
	// DeleteGroupByID deletes the group with the given id.
	DeleteGroupByID(ctx context.Context, tx pgx.Tx, groupID uuid.UUID) error
	// CreateUser adds the given store.User to the store.
	CreateUser(ctx context.Context, tx pgx.Tx, create store.User) error
	// UpdateUser updates the given store.User, identified by its id.
	UpdateUser(ctx context.Context, tx pgx.Tx, update store.User) error
	// CreateOperation creates the given store.Operation.
	CreateOperation(ctx context.Context, tx pgx.Tx, create store.Operation) error
	// UpdateOperation updates the given store.Operation.
	UpdateOperation(ctx context.Context, tx pgx.Tx, update store.Operation) error
	// UpdateOperationMembersByOperation updates the operation members for the given
	// operation.
	UpdateOperationMembersByOperation(ctx context.Context, tx pgx.Tx, operationID uuid.UUID, newMembers []uuid.UUID) error
	// DeleteForwardToGroupChannelsByGroup deletes all channels with channel type
	// store.ChannelTypeForwardToGroup, that forward to the group with the given id.
	// It returns the list of affected address book entries.
	DeleteForwardToGroupChannelsByGroup(ctx context.Context, tx pgx.Tx, groupID uuid.UUID) ([]uuid.UUID, error)
	// DeleteForwardToUserChannelsByUser deletes all channels with channel type
	// store.ChannelTypeForwardToUser, that forward to the user with the given id.
	// It returns the list of affected address book entries.
	DeleteForwardToUserChannelsByUser(ctx context.Context, tx pgx.Tx, userID uuid.UUID) ([]uuid.UUID, error)
	// CreateIntel creates the given store.Intel in the store.
	CreateIntel(ctx context.Context, tx pgx.Tx, create store.CreateIntel) (store.Intel, error)
	// IntelByID retrieves the store.Intel with the given id.
	IntelByID(ctx context.Context, tx pgx.Tx, intelID uuid.UUID) (store.Intel, error)
	// CreateIntelDelivery creates the given store.IntelDelivery in the store.
	CreateIntelDelivery(ctx context.Context, tx pgx.Tx, create store.IntelDelivery) (store.IntelDelivery, error)
	// IntelDeliveryByID retrieves the store.IntelDelivery with the given id.
	IntelDeliveryByID(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID) (store.IntelDelivery, error)
	// IntelDeliveriesTo retrieves all intel deliveries that are addressed to the address book entry entryID
	IntelDeliveriesTo(ctx context.Context, tx pgx.Tx, entryID uuid.UUID) ([]store.IntelDelivery, error)
	// TimedOutIntelDeliveryAttemptsByDelivery retrieves a
	// store.IntelDeliveryAttempt list with entries, that have been timed out (based
	// on the associated channel).
	TimedOutIntelDeliveryAttemptsByDelivery(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID) ([]store.IntelDeliveryAttempt, error)
	// UpdateIntelDeliveryAttemptStatusByID updates the status of the intel-delivery
	// attempt with the given id.
	UpdateIntelDeliveryAttemptStatusByID(ctx context.Context, tx pgx.Tx, attemptID uuid.UUID, newIsActive bool,
		newStatus store.IntelDeliveryStatus, newNote nulls.String) error
	// IntelDeliveryAttemptByID retrieves the store.IntelDeliveryAttempt with the
	// given id.
	IntelDeliveryAttemptByID(ctx context.Context, tx pgx.Tx, attemptID uuid.UUID) (store.IntelDeliveryAttempt, error)
	// NextChannelForDeliveryAttempt retrieves the next channel to try with for the
	// delivery with the given id. If none was found, the second return value will
	// be false.
	NextChannelForDeliveryAttempt(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID) (store.Channel, bool, error)
	// UpdateIntelDeliveryStatusByDelivery updates the status of the intel delivery
	// with the given id.
	UpdateIntelDeliveryStatusByDelivery(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID, newIsActive bool,
		newSuccess bool, newNote nulls.String) error
	// ActiveIntelDeliveryAttemptsByDelivery retrieves a store.IntelDeliveryAttempt
	// list with active delivery attempts.
	ActiveIntelDeliveryAttemptsByDelivery(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID) ([]store.IntelDeliveryAttempt, error)
	// CreateIntelDeliveryAttempt creates the given store.IntelDeliveryAttempt and
	// returns it with its assigned id.
	CreateIntelDeliveryAttempt(ctx context.Context, tx pgx.Tx, create store.IntelDeliveryAttempt) (store.IntelDeliveryAttempt, error)
	// LockIntelDeliveryByIDOrSkip selects the delivery with the given id in the
	// database and locks it. Selection skips locked entries, so if the entry is not
	// found or already locked, a meh.ErrNotFound will be returned.
	LockIntelDeliveryByIDOrSkip(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID) error
	// ChannelMetadataByID retrieves a store.Channel by its id without details.
	ChannelMetadataByID(ctx context.Context, tx pgx.Tx, channelID uuid.UUID) (store.Channel, error)
	// ActiveIntelDeliveryAttemptsByChannelsAndLockOrWait retrieves a
	// store.IntelDeliveryAttempt list where each one is active and uses one of the
	// given channels. It locks the associated deliveries as well as the attempts or
	// waits until locked.
	ActiveIntelDeliveryAttemptsByChannelsAndLockOrWait(ctx context.Context, tx pgx.Tx, channelIDs []uuid.UUID) ([]store.IntelDeliveryAttempt, error)
	// DeleteIntelDeliveryAttemptsByChannel deletes all intel-delivery-attempts
	// using the channel with the given id.
	DeleteIntelDeliveryAttemptsByChannel(ctx context.Context, tx pgx.Tx, channelID uuid.UUID) error
	// DeleteInactiveIntelDeliveriesFor deletes all inactive intel deliveries for a given address book entry
	DeleteInactiveIntelDeliveriesFor(ctx context.Context, tx pgx.Tx, entryID uuid.UUID) error
	// LockIntelDeliveryByIDOrWait locks the intel-delivery in the database with the
	// given id or waits until it is available.
	LockIntelDeliveryByIDOrWait(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID) error
	// ActiveIntelDeliveriesAndLockOrSkip retrieves all active intel-deliveries and
	// locks or skips them.
	ActiveIntelDeliveriesAndLockOrSkip(ctx context.Context, tx pgx.Tx) ([]store.IntelDelivery, error)
	// InvalidateIntelByID sets the valid-field of the intel with the given id to
	// false.
	InvalidateIntelByID(ctx context.Context, tx pgx.Tx, intelID uuid.UUID) error
	// SearchAddressBookEntries with the given AddressBookEntryFilters and
	// search.Params.
	SearchAddressBookEntries(ctx context.Context, tx pgx.Tx, filters store.AddressBookEntryFilters,
		searchParams search.Params) (search.Result[store.AddressBookEntryDetailed], error)
	// RebuildAddressBookEntrySearch rebuilds the address-book-entry-search.
	RebuildAddressBookEntrySearch(ctx context.Context, tx pgx.Tx) error
	// IntelDeliveryByIDAndLockOrWait retrieves the store.IntelDelivery with the
	// given id and locks it or waits until it is available.
	IntelDeliveryByIDAndLockOrWait(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID) (store.IntelDelivery, error)
	// SearchIntel using the given store.IntelFilters and search.Params.
	SearchIntel(ctx context.Context, tx pgx.Tx, filters store.IntelFilters, searchParams search.Params) (search.Result[store.Intel], error)
	// IsUserOperationMember checks if the user with the given id is member of the
	// give operation.
	IsUserOperationMember(ctx context.Context, tx pgx.Tx, userID uuid.UUID, operationID uuid.UUID) (bool, error)
	// RebuildIntelSearch rebuilds the intel-search.
	RebuildIntelSearch(ctx context.Context, tx pgx.Tx) error
	// UsersWithDeliveriesByIntel retrieves all users having associated address book
	// entries with deliveries for the intel with the given id.
	UsersWithDeliveriesByIntel(ctx context.Context, tx pgx.Tx, intelID uuid.UUID) ([]uuid.UUID, error)
	// Intel retrieves a paginated store.Intel list using the given store.IntelFilters and
	// pagination.Params, sorted descending by creation date.
	//
	// Warning: Sorting via pagination.Params is discarded!
	Intel(ctx context.Context, tx pgx.Tx, filters store.IntelFilters, paginationParams pagination.Params) (pagination.Paginated[store.Intel], error)
	// IsAutoDeliveryEnabledForAddressBookEntry checks whether the address book entry
	// with the given id is marked for auto-delivery.
	IsAutoDeliveryEnabledForAddressBookEntry(ctx context.Context, tx pgx.Tx, entryID uuid.UUID) (bool, error)
	// SetAddressBookEntriesWithAutoDeliveryEnabled sets the list of address book
	// entries with auto-delivery being enabled to the given ones. The returned list
	// of ids are of address book entries that had previously auto delivery enabled,
	// but now disabled.
	SetAddressBookEntriesWithAutoDeliveryEnabled(ctx context.Context, tx pgx.Tx, entryIDs []uuid.UUID) ([]uuid.UUID, error)
	// IntelDeliveryAttemptsByDelivery retrieves an IntelDeliveryAttempt list with
	// attempts for the delivery with the given id.
	IntelDeliveryAttemptsByDelivery(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID) ([]store.IntelDeliveryAttempt, error)
	// IntelDeliveryAttempts retrieves a paginated IntelDeliveryAttempt list using
	// the given IntelDeliveryAttemptFilters and pagination.Params, sorted descending
	// by creation date.
	//
	// Warning: Sorting via pagination.Params is discarded!
	IntelDeliveryAttempts(ctx context.Context, tx pgx.Tx, filters store.IntelDeliveryAttemptFilters,
		page pagination.Params) (pagination.Paginated[store.IntelDeliveryAttempt], error)
	// SetAutoDeliveryEnabledForAddressBookEntry sets auto intel delivery enabled for
	// the address book entry with the given id.
	SetAutoDeliveryEnabledForAddressBookEntry(ctx context.Context, tx pgx.Tx, entryID uuid.UUID, enabled bool) error
}

// Notifier sends event messages.
type Notifier interface {
	// NotifyAddressBookEntryCreated notifies created address book entries.
	NotifyAddressBookEntryCreated(ctx context.Context, tx pgx.Tx, entry store.AddressBookEntry) error
	// NotifyAddressBookEntryUpdated notifies about updated address book entries.
	NotifyAddressBookEntryUpdated(ctx context.Context, tx pgx.Tx, entry store.AddressBookEntry) error
	// NotifyAddressBookEntryDeleted notifies about deleted address book entries.
	NotifyAddressBookEntryDeleted(ctx context.Context, tx pgx.Tx, entryID uuid.UUID) error
	// NotifyAddressBookEntryChannelsUpdated notifies about updated channels for an
	// address book entry.
	NotifyAddressBookEntryChannelsUpdated(ctx context.Context, tx pgx.Tx, entryID uuid.UUID, channels []store.Channel) error
	// NotifyIntelCreated notifies about created intel.
	NotifyIntelCreated(ctx context.Context, tx pgx.Tx, created store.Intel) error
	// NotifyIntelInvalidated notifies about existing intel being invalidated.
	NotifyIntelInvalidated(ctx context.Context, tx pgx.Tx, intelID uuid.UUID, by uuid.UUID) error
	// NotifyIntelDeliveryCreated notifies about a created intel-delivery.
	NotifyIntelDeliveryCreated(ctx context.Context, tx pgx.Tx, created store.IntelDelivery) error
	// NotifyIntelDeliveryAttemptCreated notifies about a created
	// intel-delivery-attempt.
	NotifyIntelDeliveryAttemptCreated(ctx context.Context, tx pgx.Tx, created store.IntelDeliveryAttempt, delivery store.IntelDelivery,
		assignedEntry store.AddressBookEntryDetailed, intel store.Intel) error
	// NotifyIntelDeliveryAttemptStatusUpdated notifies about an status-update for a
	// intel-delivery-attempt.
	NotifyIntelDeliveryAttemptStatusUpdated(ctx context.Context, tx pgx.Tx, attempt store.IntelDeliveryAttempt) error
	// NotifyIntelDeliveryStatusUpdated notifies abous a status-update for an
	// intel-delivery.
	NotifyIntelDeliveryStatusUpdated(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID, newIsActive bool,
		newSuccess bool, newNote nulls.String) error
	// NotifyAddressBookEntryAutoDeliveryUpdated emits an
	// event.TypeAddressBookEntryAutoDeliveryUpdated event.
	NotifyAddressBookEntryAutoDeliveryUpdated(ctx context.Context, tx pgx.Tx, entryID uuid.UUID, isAutoDeliveryEnabled bool) error
}
