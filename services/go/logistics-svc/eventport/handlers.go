package eventport

import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/lefinal/meh"
	"github.com/lefinal/nulls"
	"github.com/mobile-directing-system/mds-server/services/go/logistics-svc/store"
	"github.com/mobile-directing-system/mds-server/services/go/shared/event"
	"github.com/mobile-directing-system/mds-server/services/go/shared/kafkautil"
)

// Handler for received messages.
type Handler interface {
	// CreateUser creates the user with the given id.
	CreateUser(ctx context.Context, tx pgx.Tx, userID store.User) error
	// UpdateUser updates the given store.user, identified by its id.
	UpdateUser(ctx context.Context, tx pgx.Tx, user store.User) error
	// CreateGroup creates the given store.Group.
	CreateGroup(ctx context.Context, tx pgx.Tx, create store.Group) error
	// UpdateGroup updates the given store.Group, identified by its id.
	UpdateGroup(ctx context.Context, tx pgx.Tx, update store.Group) error
	// DeleteGroupByID deletes the group with the given id.
	DeleteGroupByID(ctx context.Context, tx pgx.Tx, groupID uuid.UUID) error
	// CreateOperation creates the given store.Operation.
	CreateOperation(ctx context.Context, tx pgx.Tx, create store.Operation) error
	// UpdateOperation updates the given store.Operation.
	UpdateOperation(ctx context.Context, tx pgx.Tx, update store.Operation) error
	// UpdateOperationMembersByOperation updates the operation members for the given
	// operation.
	UpdateOperationMembersByOperation(ctx context.Context, tx pgx.Tx, operationID uuid.UUID, newMembers []uuid.UUID) error
	// CreateIntel and schedule delivery.
	CreateIntel(ctx context.Context, tx pgx.Tx, create store.Intel) error
	// InvalidateIntelByID sets the valid-field of the intel with the given id to
	// false.
	InvalidateIntelByID(ctx context.Context, tx pgx.Tx, intelID uuid.UUID) error
	// UpdateIntelDeliveryAttemptStatusForActive updates the
	// intel-delivery-attempt-status for the attempt with the given id. It assures
	// that the delivery attempt is still active and does not have
	// store.IntelDeliveryStatusCanceled.
	UpdateIntelDeliveryAttemptStatusForActive(ctx context.Context, tx pgx.Tx, attemptID uuid.UUID,
		newStatus store.IntelDeliveryStatus, newNote nulls.String) error
}

// HandlerFn for handling messages.
func (p *Port) HandlerFn(handler Handler) kafkautil.HandlerFunc {
	return func(ctx context.Context, tx pgx.Tx, message kafkautil.InboundMessage) error {
		switch message.Topic {
		case event.GroupsTopic:
			return meh.NilOrWrap(p.handleGroupsTopic(ctx, tx, handler, message), "handle groups topic", nil)
		case event.OperationsTopic:
			return meh.NilOrWrap(p.handleOperationsTopic(ctx, tx, handler, message), "handle operations topic", nil)
		case event.UsersTopic:
			return meh.NilOrWrap(p.handleUsersTopic(ctx, tx, handler, message), "handle users topic", nil)
		case event.IntelTopic:
			return meh.NilOrWrap(p.handleIntelTopic(ctx, tx, handler, message), "handle intel topic", nil)
		case event.InAppNotificationsTopic:
			return meh.NilOrWrap(p.handleInAppNotificationsTopic(ctx, tx, handler, message), "handle in-app-notifications topic", nil)
		}
		return nil
	}
}

// handleGroupsTopic handles the event.GroupsTopic.
func (p *Port) handleGroupsTopic(ctx context.Context, tx pgx.Tx, handler Handler, message kafkautil.InboundMessage) error {
	switch message.EventType {
	case event.TypeGroupCreated:
		return meh.NilOrWrap(p.handleGroupCreated(ctx, tx, handler, message), "handle group created", nil)
	case event.TypeGroupDeleted:
		return meh.NilOrWrap(p.handleGroupDeleted(ctx, tx, handler, message), "handle group deleted", nil)
	case event.TypeGroupUpdated:
		return meh.NilOrWrap(p.handleGroupUpdated(ctx, tx, handler, message), "handle group updated", nil)
	}
	return nil
}

// handleGroupCreated handles an event.TypeGroupCreated event.
func (p *Port) handleGroupCreated(ctx context.Context, tx pgx.Tx, handler Handler, message kafkautil.InboundMessage) error {
	var groupCreatedEvent event.GroupCreated
	err := json.Unmarshal(message.RawValue, &groupCreatedEvent)
	if err != nil {
		return meh.NewInternalErrFromErr(err, "unmarshal event", meh.Details{"raw": string(message.RawValue)})
	}
	create := store.Group{
		ID:          groupCreatedEvent.ID,
		Title:       groupCreatedEvent.Title,
		Description: groupCreatedEvent.Description,
		Operation:   groupCreatedEvent.Operation,
		Members:     groupCreatedEvent.Members,
	}
	err = handler.CreateGroup(ctx, tx, create)
	if err != nil {
		return meh.Wrap(err, "create group", meh.Details{"group": create})
	}
	return nil
}

// handleGroupUpdated handles an event.TypeGroupUpdated event.
func (p *Port) handleGroupUpdated(ctx context.Context, tx pgx.Tx, handler Handler, message kafkautil.InboundMessage) error {
	var groupUpdatedEvent event.GroupUpdated
	err := json.Unmarshal(message.RawValue, &groupUpdatedEvent)
	if err != nil {
		return meh.NewInternalErrFromErr(err, "unmarshal event", meh.Details{"raw": string(message.RawValue)})
	}
	update := store.Group{
		ID:          groupUpdatedEvent.ID,
		Title:       groupUpdatedEvent.Title,
		Description: groupUpdatedEvent.Description,
		Operation:   groupUpdatedEvent.Operation,
		Members:     groupUpdatedEvent.Members,
	}
	err = handler.UpdateGroup(ctx, tx, update)
	if err != nil {
		return meh.Wrap(err, "update group", meh.Details{"group": update})
	}
	return nil
}

// handleGroupDeleted handles an event.TypeGroupDeleted event.
func (p *Port) handleGroupDeleted(ctx context.Context, tx pgx.Tx, handler Handler, message kafkautil.InboundMessage) error {
	var groupDeletedEvent event.GroupDeleted
	err := json.Unmarshal(message.RawValue, &groupDeletedEvent)
	if err != nil {
		return meh.NewInternalErrFromErr(err, "unmarshal event", meh.Details{"raw": string(message.RawValue)})
	}
	err = handler.DeleteGroupByID(ctx, tx, groupDeletedEvent.ID)
	if err != nil {
		return meh.Wrap(err, "delete group", meh.Details{"group_id": groupDeletedEvent.ID})
	}
	return nil
}

// handleUsersTopic handles the event.UsersTopic.
func (p *Port) handleUsersTopic(ctx context.Context, tx pgx.Tx, handler Handler, message kafkautil.InboundMessage) error {
	switch message.EventType {
	case event.TypeUserCreated:
		return meh.NilOrWrap(p.handleUserCreated(ctx, tx, handler, message), "handle user created", nil)
	case event.TypeUserUpdated:
		return meh.NilOrWrap(p.handleUserUpdated(ctx, tx, handler, message), "handle user updated", nil)
	}
	return nil
}

// handleUserCreated handles an event.TypeUserCreated event.
func (p *Port) handleUserCreated(ctx context.Context, tx pgx.Tx, handler Handler, message kafkautil.InboundMessage) error {
	var userCreatedEvent event.UserCreated
	err := json.Unmarshal(message.RawValue, &userCreatedEvent)
	if err != nil {
		return meh.NewInternalErrFromErr(err, "unmarshal event", meh.Details{"raw": string(message.RawValue)})
	}
	create := store.User{
		ID:        userCreatedEvent.ID,
		Username:  userCreatedEvent.Username,
		FirstName: userCreatedEvent.FirstName,
		LastName:  userCreatedEvent.LastName,
		IsActive:  userCreatedEvent.IsActive,
	}
	err = handler.CreateUser(ctx, tx, create)
	if err != nil {
		return meh.Wrap(err, "create user", meh.Details{"user": create})
	}
	return nil
}

// handleUserUpdated handles an event.TypeUserUpdated event.
func (p *Port) handleUserUpdated(ctx context.Context, tx pgx.Tx, handler Handler, message kafkautil.InboundMessage) error {
	var userUpdatedEvent event.UserUpdated
	err := json.Unmarshal(message.RawValue, &userUpdatedEvent)
	if err != nil {
		return meh.NewInternalErrFromErr(err, "unmarshal event", meh.Details{"raw": string(message.RawValue)})
	}
	update := store.User{
		ID:        userUpdatedEvent.ID,
		Username:  userUpdatedEvent.Username,
		FirstName: userUpdatedEvent.FirstName,
		LastName:  userUpdatedEvent.LastName,
		IsActive:  userUpdatedEvent.IsActive,
	}
	err = handler.UpdateUser(ctx, tx, update)
	if err != nil {
		return meh.Wrap(err, "update user", meh.Details{"user": update})
	}
	return nil
}

// handleOperationsTopic handles the event.OperationsTopic.
func (p *Port) handleOperationsTopic(ctx context.Context, tx pgx.Tx, handler Handler, message kafkautil.InboundMessage) error {
	switch message.EventType {
	case event.TypeOperationCreated:
		return meh.NilOrWrap(p.handleOperationCreated(ctx, tx, handler, message), "handle operation created", nil)
	case event.TypeOperationUpdated:
		return meh.NilOrWrap(p.handleOperationUpdated(ctx, tx, handler, message), "handle operation updated", nil)
	case event.TypeOperationMembersUpdated:
		return meh.NilOrWrap(p.handleOperationMembersUpdated(ctx, tx, handler, message), "handle operation members updated", nil)
	}
	return nil
}

// handleOperationCreated handles an event.TypeOperationCreated event.
func (p *Port) handleOperationCreated(ctx context.Context, tx pgx.Tx, handler Handler, message kafkautil.InboundMessage) error {
	var operationCreatedEvent event.OperationCreated
	err := json.Unmarshal(message.RawValue, &operationCreatedEvent)
	if err != nil {
		return meh.NewInternalErrFromErr(err, "unmarshal event", meh.Details{"raw": string(message.RawValue)})
	}
	create := store.Operation{
		ID:          operationCreatedEvent.ID,
		Title:       operationCreatedEvent.Title,
		Description: operationCreatedEvent.Description,
		Start:       operationCreatedEvent.Start,
		End:         operationCreatedEvent.End,
		IsArchived:  operationCreatedEvent.IsArchived,
	}
	err = handler.CreateOperation(ctx, tx, create)
	if err != nil {
		return meh.Wrap(err, "create operation", meh.Details{"create": create})
	}
	return nil
}

// handleOperationUpdated handles an event.TypeOperationUpdated event.
func (p *Port) handleOperationUpdated(ctx context.Context, tx pgx.Tx, handler Handler, message kafkautil.InboundMessage) error {
	var operationUpdatedEvent event.OperationUpdated
	err := json.Unmarshal(message.RawValue, &operationUpdatedEvent)
	if err != nil {
		return meh.NewInternalErrFromErr(err, "unmarshal event", meh.Details{"raw": string(message.RawValue)})
	}
	update := store.Operation{
		ID:          operationUpdatedEvent.ID,
		Title:       operationUpdatedEvent.Title,
		Description: operationUpdatedEvent.Description,
		Start:       operationUpdatedEvent.Start,
		End:         operationUpdatedEvent.End,
		IsArchived:  operationUpdatedEvent.IsArchived,
	}
	err = handler.UpdateOperation(ctx, tx, update)
	if err != nil {
		return meh.Wrap(err, "update operation", meh.Details{"update": update})
	}
	return nil
}

// handleOperationMembersUpdated handles an event.TypeOperationMembersUpdated
// event.
func (p *Port) handleOperationMembersUpdated(ctx context.Context, tx pgx.Tx, handler Handler, message kafkautil.InboundMessage) error {
	var operationMembersUpdatedEvent event.OperationMembersUpdated
	err := json.Unmarshal(message.RawValue, &operationMembersUpdatedEvent)
	if err != nil {
		return meh.NewInternalErrFromErr(err, "unmarshal event", meh.Details{"raw": string(message.RawValue)})
	}
	err = handler.UpdateOperationMembersByOperation(ctx, tx, operationMembersUpdatedEvent.Operation, operationMembersUpdatedEvent.Members)
	if err != nil {
		return meh.Wrap(err, "update operation members", meh.Details{
			"operation":   operationMembersUpdatedEvent.Operation,
			"new_members": operationMembersUpdatedEvent.Members,
		})
	}
	return nil
}

// handleIntelTopic handles the event.IntelTopic.
func (p *Port) handleIntelTopic(ctx context.Context, tx pgx.Tx, handler Handler, message kafkautil.InboundMessage) error {
	switch message.EventType {
	case event.TypeIntelCreated:
		return meh.NilOrWrap(p.handleIntelCreated(ctx, tx, handler, message), "handle intel created", nil)
	case event.TypeIntelInvalidated:
		return meh.NilOrWrap(p.handleIntelInvalidated(ctx, tx, handler, message), "handle intel invalidated", nil)
	}
	return nil
}

// handleIntelCreated handles an event.TypeIntelCreated event.
func (p *Port) handleIntelCreated(ctx context.Context, tx pgx.Tx, handler Handler, message kafkautil.InboundMessage) error {
	var intelCreatedEvent event.IntelCreated
	err := json.Unmarshal(message.RawValue, &intelCreatedEvent)
	if err != nil {
		return meh.NewInternalErrFromErr(err, "unmarshal event", meh.Details{"raw": string(message.RawValue)})
	}
	sAssignments := make([]store.IntelAssignment, 0, len(intelCreatedEvent.Assignments))
	for _, assignment := range intelCreatedEvent.Assignments {
		sAssignments = append(sAssignments, store.IntelAssignment{
			ID:    assignment.ID,
			Intel: intelCreatedEvent.ID,
			To:    assignment.To,
		})
	}
	sCreate := store.Intel{
		ID:          intelCreatedEvent.ID,
		CreatedAt:   intelCreatedEvent.CreatedAt,
		CreatedBy:   intelCreatedEvent.CreatedBy,
		Operation:   intelCreatedEvent.Operation,
		Type:        store.IntelType(intelCreatedEvent.Type), // No conversation, as we simply do not care.
		Content:     intelCreatedEvent.Content,
		SearchText:  intelCreatedEvent.SearchText,
		Importance:  intelCreatedEvent.Importance,
		IsValid:     intelCreatedEvent.IsValid,
		Assignments: sAssignments,
	}
	err = handler.CreateIntel(ctx, tx, sCreate)
	if err != nil {
		return meh.Wrap(err, "create intel", meh.Details{"create": sCreate})
	}
	return nil
}

// handleIntelInvalidated handles an event.TypeIntelInvalidated event.
func (p *Port) handleIntelInvalidated(ctx context.Context, tx pgx.Tx, handler Handler, message kafkautil.InboundMessage) error {
	var intelInvalidatedEvent event.IntelInvalidated
	err := json.Unmarshal(message.RawValue, &intelInvalidatedEvent)
	if err != nil {
		return meh.NewInternalErrFromErr(err, "unmarshal event", meh.Details{"raw": string(message.RawValue)})
	}
	err = handler.InvalidateIntelByID(ctx, tx, intelInvalidatedEvent.ID)
	if err != nil {
		return meh.Wrap(err, "create intel", meh.Details{"intel_id": intelInvalidatedEvent.ID})
	}
	return nil
}

// handleInAppNotificationsTopic handles the event.InAppNotificationsTopic.
func (p *Port) handleInAppNotificationsTopic(ctx context.Context, tx pgx.Tx, handler Handler, message kafkautil.InboundMessage) error {
	switch message.EventType {
	case event.TypeInAppNotificationForIntelPending:
		return meh.NilOrWrap(p.handleInAppNotificationForIntelPending(ctx, tx, handler, message), "handle in-app-notification pending", nil)
	case event.TypeInAppNotificationForIntelSent:
		return meh.NilOrWrap(p.handleInAppNotificationForIntelSent(ctx, tx, handler, message), "handle in-app-notification sent", nil)
	}
	return nil
}

// handleInAppNotificationForIntelPending handles an
// event.TypeInAppNotificationForIntelPending event.
func (p *Port) handleInAppNotificationForIntelPending(ctx context.Context, tx pgx.Tx, handler Handler, message kafkautil.InboundMessage) error {
	var notifPendingEvent event.InAppNotificationForIntelPending
	err := json.Unmarshal(message.RawValue, &notifPendingEvent)
	if err != nil {
		return meh.NewInternalErrFromErr(err, "unmarshal event", meh.Details{"raw": string(message.RawValue)})
	}
	err = handler.UpdateIntelDeliveryAttemptStatusForActive(ctx, tx, notifPendingEvent.Attempt, store.IntelDeliveryStatusAwaitingDelivery,
		nulls.NewString("in-app-notification pending"))
	if err != nil {
		return meh.Wrap(err, "update intel-delivery-attempt-status for active", meh.Details{"attempt_id": notifPendingEvent.Attempt})
	}
	return nil
}

// handleInAppNotificationForIntelSent handles an
// event.TypeInAppNotificationForIntelSent event.
func (p *Port) handleInAppNotificationForIntelSent(ctx context.Context, tx pgx.Tx, handler Handler, message kafkautil.InboundMessage) error {
	var notifSentEvent event.InAppNotificationForIntelSent
	err := json.Unmarshal(message.RawValue, &notifSentEvent)
	if err != nil {
		return meh.NewInternalErrFromErr(err, "unmarshal event", meh.Details{"raw": string(message.RawValue)})
	}
	err = handler.UpdateIntelDeliveryAttemptStatusForActive(ctx, tx, notifSentEvent.Attempt, store.IntelDeliveryStatusAwaitingAck,
		nulls.NewString("in-app-notification sent"))
	if err != nil {
		return meh.Wrap(err, "update intel-delivery-attempt-status for active", meh.Details{"attempt_id": notifSentEvent.Attempt})
	}
	return nil
}