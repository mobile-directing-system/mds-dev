package controller

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/lefinal/nulls"
	"github.com/mobile-directing-system/mds-server/services/go/logistics-svc/store"
	"github.com/mobile-directing-system/mds-server/services/go/shared/pagination"
	"github.com/mobile-directing-system/mds-server/services/go/shared/search"
	"github.com/mobile-directing-system/mds-server/services/go/shared/testutil"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"time"
)

const timeout = 5 * time.Second

type ControllerMock struct {
	Logger   *zap.Logger
	DB       *testutil.DBTxSupplier
	Store    *StoreMock
	Notifier *NotifierMock
	Ctrl     *Controller
}

func NewMockController() *ControllerMock {
	ctrl := &ControllerMock{
		Logger:   zap.NewNop(),
		DB:       &testutil.DBTxSupplier{},
		Store:    &StoreMock{},
		Notifier: &NotifierMock{},
	}
	ctrl.Ctrl = &Controller{
		Logger:   ctrl.Logger,
		DB:       ctrl.DB,
		Store:    ctrl.Store,
		Notifier: ctrl.Notifier,
	}
	return ctrl
}

// StoreMock mocks Store.
type StoreMock struct {
	mock.Mock
}

func (m *StoreMock) ChannelsByAddressBookEntry(ctx context.Context, tx pgx.Tx, entryID uuid.UUID) ([]store.Channel, error) {
	args := m.Called(ctx, tx, entryID)
	var channels []store.Channel
	channels, _ = args.Get(0).([]store.Channel)
	return channels, args.Error(1)
}

func (m *StoreMock) AssureAddressBookEntryExists(ctx context.Context, tx pgx.Tx, entryID uuid.UUID) error {
	return m.Called(ctx, tx, entryID).Error(0)
}

func (m *StoreMock) DeleteChannelWithDetailsByID(ctx context.Context, tx pgx.Tx, channelID uuid.UUID, channelType store.ChannelType) error {
	return m.Called(ctx, tx, channelID, channelType).Error(0)
}

func (m *StoreMock) AddressBookEntries(ctx context.Context, tx pgx.Tx, filters store.AddressBookEntryFilters,
	paginationParams pagination.Params) (pagination.Paginated[store.AddressBookEntryDetailed], error) {
	args := m.Called(ctx, tx, filters, paginationParams)
	return args.Get(0).(pagination.Paginated[store.AddressBookEntryDetailed]), args.Error(1)
}

func (m *StoreMock) CreateAddressBookEntry(ctx context.Context, tx pgx.Tx, entry store.AddressBookEntry) (store.AddressBookEntryDetailed, error) {
	args := m.Called(ctx, tx, entry)
	return args.Get(0).(store.AddressBookEntryDetailed), args.Error(1)
}

func (m *StoreMock) UpdateAddressBookEntry(ctx context.Context, tx pgx.Tx, entry store.AddressBookEntry) error {
	return m.Called(ctx, tx, entry).Error(0)
}

func (m *StoreMock) DeleteAddressBookEntryByID(ctx context.Context, tx pgx.Tx, entryID uuid.UUID) error {
	return m.Called(ctx, tx, entryID).Error(0)
}

func (m *StoreMock) CreateGroup(ctx context.Context, tx pgx.Tx, create store.Group) error {
	return m.Called(ctx, tx, create).Error(0)
}

func (m *StoreMock) UpdateGroup(ctx context.Context, tx pgx.Tx, update store.Group) error {
	return m.Called(ctx, tx, update).Error(0)
}

func (m *StoreMock) DeleteGroupByID(ctx context.Context, tx pgx.Tx, groupID uuid.UUID) error {
	return m.Called(ctx, tx, groupID).Error(0)
}

func (m *StoreMock) CreateUser(ctx context.Context, tx pgx.Tx, create store.User) error {
	return m.Called(ctx, tx, create).Error(0)
}

func (m *StoreMock) UpdateUser(ctx context.Context, tx pgx.Tx, update store.User) error {
	return m.Called(ctx, tx, update).Error(0)
}

func (m *StoreMock) DeleteUserByID(ctx context.Context, tx pgx.Tx, userID uuid.UUID) error {
	return m.Called(ctx, tx, userID).Error(0)
}

func (m *StoreMock) CreateOperation(ctx context.Context, tx pgx.Tx, create store.Operation) error {
	return m.Called(ctx, tx, create).Error(0)
}

func (m *StoreMock) UpdateOperation(ctx context.Context, tx pgx.Tx, update store.Operation) error {
	return m.Called(ctx, tx, update).Error(0)
}

func (m *StoreMock) UpdateOperationMembersByOperation(ctx context.Context, tx pgx.Tx, operationID uuid.UUID, newMembers []uuid.UUID) error {
	return m.Called(ctx, tx, operationID, newMembers).Error(0)
}

func (m *StoreMock) AddressBookEntryByID(ctx context.Context, tx pgx.Tx, entryID uuid.UUID,
	visibleBy uuid.NullUUID) (store.AddressBookEntryDetailed, error) {
	args := m.Called(ctx, tx, entryID, visibleBy)
	return args.Get(0).(store.AddressBookEntryDetailed), args.Error(1)
}

func (m *StoreMock) DeleteForwardToGroupChannelsByGroup(ctx context.Context, tx pgx.Tx, groupID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, tx, groupID)
	var channels []uuid.UUID
	channels, _ = args.Get(0).([]uuid.UUID)
	return channels, args.Error(1)
}

func (m *StoreMock) DeleteForwardToUserChannelsByUser(ctx context.Context, tx pgx.Tx, userID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, tx, userID)
	var channels []uuid.UUID
	channels, _ = args.Get(0).([]uuid.UUID)
	return channels, args.Error(1)
}

func (m *StoreMock) UpdateChannelsByEntry(ctx context.Context, tx pgx.Tx, entryID uuid.UUID, newChannels []store.Channel) error {
	return m.Called(ctx, tx, entryID, newChannels).Error(0)
}

func (m *StoreMock) IntelByID(ctx context.Context, tx pgx.Tx, intelID uuid.UUID) (store.Intel, error) {
	args := m.Called(ctx, tx, intelID)
	return args.Get(0).(store.Intel), args.Error(1)
}

func (m *StoreMock) CreateIntelDelivery(ctx context.Context, tx pgx.Tx, create store.IntelDelivery) (store.IntelDelivery, error) {
	args := m.Called(ctx, tx, create)
	return args.Get(0).(store.IntelDelivery), args.Error(1)
}

func (m *StoreMock) IntelDeliveryByID(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID) (store.IntelDelivery, error) {
	args := m.Called(ctx, tx, deliveryID)
	return args.Get(0).(store.IntelDelivery), args.Error(1)
}

func (m *StoreMock) IntelDeliveriesTo(ctx context.Context, tx pgx.Tx, entryID uuid.UUID) ([]store.IntelDelivery, error) {
	args := m.Called(ctx, tx, entryID)
	deliveries, _ := args.Get(0).([]store.IntelDelivery)
	return deliveries, args.Error(1)
}

func (m *StoreMock) TimedOutIntelDeliveryAttemptsByDelivery(ctx context.Context, tx pgx.Tx,
	deliveryID uuid.UUID) ([]store.IntelDeliveryAttempt, error) {
	args := m.Called(ctx, tx, deliveryID)
	var attempts []store.IntelDeliveryAttempt
	if a := args.Get(0); a != nil {
		attempts = a.([]store.IntelDeliveryAttempt)
	}
	return attempts, args.Error(1)
}

func (m *StoreMock) UpdateIntelDeliveryAttemptStatusByID(ctx context.Context, tx pgx.Tx, attemptID uuid.UUID,
	newIsActive bool, newStatus store.IntelDeliveryStatus, newNote nulls.String) error {
	return m.Called(ctx, tx, attemptID, newIsActive, newStatus, newNote).Error(0)
}

func (m *StoreMock) IntelDeliveryAttemptByID(ctx context.Context, tx pgx.Tx, attemptID uuid.UUID) (store.IntelDeliveryAttempt, error) {
	args := m.Called(ctx, tx, attemptID)
	return args.Get(0).(store.IntelDeliveryAttempt), args.Error(1)
}

func (m *StoreMock) NextChannelForDeliveryAttempt(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID) (store.Channel, bool, error) {
	args := m.Called(ctx, tx, deliveryID)
	return args.Get(0).(store.Channel), args.Bool(1), args.Error(2)
}

func (m *StoreMock) UpdateIntelDeliveryStatusByDelivery(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID,
	newIsActive bool, newSuccess bool, newNote nulls.String) error {
	return m.Called(ctx, tx, deliveryID, newIsActive, newSuccess, newNote).Error(0)
}

func (m *StoreMock) ActiveIntelDeliveryAttemptsByDelivery(ctx context.Context, tx pgx.Tx,
	deliveryID uuid.UUID) ([]store.IntelDeliveryAttempt, error) {
	args := m.Called(ctx, tx, deliveryID)
	var attempts []store.IntelDeliveryAttempt
	if a := args.Get(0); a != nil {
		attempts = a.([]store.IntelDeliveryAttempt)
	}
	return attempts, args.Error(1)
}

func (m *StoreMock) CreateIntelDeliveryAttempt(ctx context.Context, tx pgx.Tx,
	create store.IntelDeliveryAttempt) (store.IntelDeliveryAttempt, error) {
	args := m.Called(ctx, tx, create)
	return args.Get(0).(store.IntelDeliveryAttempt), args.Error(1)
}

func (m *StoreMock) LockIntelDeliveryByIDOrSkip(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID) error {
	return m.Called(ctx, tx, deliveryID).Error(0)
}

func (m *StoreMock) ChannelMetadataByID(ctx context.Context, tx pgx.Tx, channelID uuid.UUID) (store.Channel, error) {
	args := m.Called(ctx, tx, channelID)
	return args.Get(0).(store.Channel), args.Error(1)
}

func (m *StoreMock) ActiveIntelDeliveryAttemptsByChannelsAndLockOrWait(ctx context.Context, tx pgx.Tx,
	channelIDs []uuid.UUID) ([]store.IntelDeliveryAttempt, error) {
	args := m.Called(ctx, tx, channelIDs)
	var attempts []store.IntelDeliveryAttempt
	if a := args.Get(0); a != nil {
		attempts = a.([]store.IntelDeliveryAttempt)
	}
	return attempts, args.Error(1)
}

func (m *StoreMock) DeleteIntelDeliveryAttemptsByChannel(ctx context.Context, tx pgx.Tx, channelID uuid.UUID) error {
	return m.Called(ctx, tx, channelID).Error(0)
}

func (m *StoreMock) DeleteInactiveIntelDeliveriesFor(ctx context.Context, tx pgx.Tx, entryID uuid.UUID) error {
	return m.Called(ctx, tx, entryID).Error(0)
}

func (m *StoreMock) LockIntelDeliveryByIDOrWait(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID) error {
	return m.Called(ctx, tx, deliveryID).Error(0)
}

func (m *StoreMock) ActiveIntelDeliveriesAndLockOrSkip(ctx context.Context, tx pgx.Tx) ([]store.IntelDelivery, error) {
	args := m.Called(ctx, tx)
	var deliveries []store.IntelDelivery
	if a := args.Get(0); a != nil {
		deliveries = a.([]store.IntelDelivery)
	}
	return deliveries, args.Error(1)
}

func (m *StoreMock) InvalidateIntelByID(ctx context.Context, tx pgx.Tx, intelID uuid.UUID) error {
	return m.Called(ctx, tx, intelID).Error(0)
}

func (m *StoreMock) SearchAddressBookEntries(ctx context.Context, tx pgx.Tx, filters store.AddressBookEntryFilters,
	searchParams search.Params) (search.Result[store.AddressBookEntryDetailed], error) {
	args := m.Called(ctx, tx, filters, searchParams)
	return args.Get(0).(search.Result[store.AddressBookEntryDetailed]), args.Error(1)
}

func (m *StoreMock) RebuildAddressBookEntrySearch(ctx context.Context, tx pgx.Tx) error {
	return m.Called(ctx, tx).Error(0)
}

func (m *StoreMock) IntelDeliveryByIDAndLockOrWait(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID) (store.IntelDelivery, error) {
	args := m.Called(ctx, tx, deliveryID)
	return args.Get(0).(store.IntelDelivery), args.Error(1)
}

func (m *StoreMock) CreateIntel(ctx context.Context, tx pgx.Tx, create store.CreateIntel) (store.Intel, error) {
	args := m.Called(ctx, tx, create)
	return args.Get(0).(store.Intel), args.Error(1)
}

func (m *StoreMock) SearchIntel(ctx context.Context, tx pgx.Tx, filters store.IntelFilters,
	searchParams search.Params) (search.Result[store.Intel], error) {
	args := m.Called(ctx, tx, filters, searchParams)
	return args.Get(0).(search.Result[store.Intel]), args.Error(1)
}

func (m *StoreMock) IsUserOperationMember(ctx context.Context, tx pgx.Tx, userID uuid.UUID, operationID uuid.UUID) (bool, error) {
	args := m.Called(ctx, tx, userID, operationID)
	return args.Bool(0), args.Error(1)
}

func (m *StoreMock) RebuildIntelSearch(ctx context.Context, tx pgx.Tx) error {
	return m.Called(ctx, tx).Error(0)
}

func (m *StoreMock) UsersWithDeliveriesByIntel(ctx context.Context, tx pgx.Tx, intelID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, tx, intelID)
	var users []uuid.UUID
	if v := args.Get(0); v != nil {
		users = v.([]uuid.UUID)
	}
	return users, args.Error(1)
}

func (m *StoreMock) Intel(ctx context.Context, tx pgx.Tx, filters store.IntelFilters,
	paginationParams pagination.Params) (pagination.Paginated[store.Intel], error) {
	args := m.Called(ctx, tx, filters, paginationParams)
	return args.Get(0).(pagination.Paginated[store.Intel]), args.Error(1)
}

func (m *StoreMock) IsAutoDeliveryEnabledForAddressBookEntry(ctx context.Context, tx pgx.Tx, entryID uuid.UUID) (bool, error) {
	args := m.Called(ctx, tx, entryID)
	return args.Bool(0), args.Error(1)
}

func (m *StoreMock) SetAddressBookEntriesWithAutoDeliveryEnabled(ctx context.Context, tx pgx.Tx, entryIDs []uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, tx, entryIDs)
	var disabled []uuid.UUID
	if a := args.Get(0); a != nil {
		disabled = a.([]uuid.UUID)
	}
	return disabled, args.Error(1)
}

func (m *StoreMock) IntelDeliveryAttemptsByDelivery(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID) ([]store.IntelDeliveryAttempt, error) {
	args := m.Called(ctx, tx, deliveryID)
	var attempts []store.IntelDeliveryAttempt
	if a := args.Get(0); a != nil {
		attempts = a.([]store.IntelDeliveryAttempt)
	}
	return attempts, args.Error(1)
}

func (m *StoreMock) IntelDeliveryAttempts(ctx context.Context, tx pgx.Tx, filters store.IntelDeliveryAttemptFilters,
	page pagination.Params) (pagination.Paginated[store.IntelDeliveryAttempt], error) {
	args := m.Called(ctx, tx, filters, page)
	return args.Get(0).(pagination.Paginated[store.IntelDeliveryAttempt]), args.Error(1)
}

func (m *StoreMock) SetAutoDeliveryEnabledForAddressBookEntry(ctx context.Context, tx pgx.Tx, entryID uuid.UUID, enabled bool) error {
	return m.Called(ctx, tx, entryID, enabled).Error(0)
}

// NotifierMock mocks Notifier.
type NotifierMock struct {
	mock.Mock
}

func (m *NotifierMock) NotifyAddressBookEntryCreated(ctx context.Context, tx pgx.Tx, entry store.AddressBookEntry) error {
	return m.Called(ctx, tx, entry).Error(0)
}

func (m *NotifierMock) NotifyAddressBookEntryUpdated(ctx context.Context, tx pgx.Tx, entry store.AddressBookEntry) error {
	return m.Called(ctx, tx, entry).Error(0)
}

func (m *NotifierMock) NotifyAddressBookEntryDeleted(ctx context.Context, tx pgx.Tx, entryID uuid.UUID) error {
	return m.Called(ctx, tx, entryID).Error(0)
}

func (m *NotifierMock) NotifyAddressBookEntryChannelsUpdated(ctx context.Context, tx pgx.Tx, entryID uuid.UUID, channels []store.Channel) error {
	return m.Called(ctx, tx, entryID, channels).Error(0)
}

func (m *NotifierMock) NotifyIntelDeliveryCreated(ctx context.Context, tx pgx.Tx, created store.IntelDelivery) error {
	return m.Called(ctx, tx, created).Error(0)
}

func (m *NotifierMock) NotifyIntelDeliveryAttemptCreated(ctx context.Context, tx pgx.Tx, created store.IntelDeliveryAttempt,
	delivery store.IntelDelivery, assignedEntry store.AddressBookEntryDetailed, intel store.Intel) error {
	return m.Called(ctx, tx, created, delivery, assignedEntry, intel).Error(0)
}

func (m *NotifierMock) NotifyIntelDeliveryAttemptStatusUpdated(ctx context.Context, tx pgx.Tx, attempt store.IntelDeliveryAttempt) error {
	return m.Called(ctx, tx, attempt).Error(0)
}

func (m *NotifierMock) NotifyIntelDeliveryStatusUpdated(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID, newIsActive bool,
	newSuccess bool, newNote nulls.String) error {
	return m.Called(ctx, tx, deliveryID, newIsActive, newSuccess, newNote).Error(0)
}

func (m *NotifierMock) NotifyIntelCreated(ctx context.Context, tx pgx.Tx, created store.Intel) error {
	return m.Called(ctx, tx, created).Error(0)
}

func (m *NotifierMock) NotifyIntelInvalidated(ctx context.Context, tx pgx.Tx, intelID uuid.UUID, by uuid.UUID) error {
	return m.Called(ctx, tx, intelID, by).Error(0)
}

func (m *NotifierMock) NotifyAddressBookEntryAutoDeliveryUpdated(ctx context.Context, tx pgx.Tx, entryID uuid.UUID, isAutoDeliveryEnabled bool) error {
	return m.Called(ctx, tx, entryID, isAutoDeliveryEnabled).Error(0)
}
