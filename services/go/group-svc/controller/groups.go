package controller

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/lefinal/meh"
	"github.com/mobile-directing-system/mds-server/services/go/group-svc/store"
	"github.com/mobile-directing-system/mds-server/services/go/shared/pagination"
	"github.com/mobile-directing-system/mds-server/services/go/shared/pgutil"
)

// CreateGroup creates and notifies of the given group.
func (c *Controller) CreateGroup(ctx context.Context, group store.Group) (store.Group, error) {
	var created store.Group
	err := pgutil.RunInTx(ctx, c.DB, func(ctx context.Context, tx pgx.Tx) error {
		// Create in store.
		var err error
		created, err = c.Store.CreateGroup(ctx, tx, group)
		if err != nil {
			return meh.Wrap(err, "create group in store", meh.Details{"group": group})
		}
		// Notify.
		err = c.Notifier.NotifyGroupCreated(created)
		if err != nil {
			return meh.Wrap(err, "notify", meh.Details{"group": created})
		}
		return nil
	})
	if err != nil {
		return store.Group{}, meh.Wrap(err, "run in tx", nil)
	}
	return created, nil
}

// UpdateGroup updates the given store.Group, identifed by its id, and notifies.
func (c *Controller) UpdateGroup(ctx context.Context, group store.Group) error {
	err := pgutil.RunInTx(ctx, c.DB, func(ctx context.Context, tx pgx.Tx) error {
		// Update in store.
		err := c.Store.UpdateGroup(ctx, tx, group)
		if err != nil {
			return meh.Wrap(err, "update group in store", meh.Details{"group": group})
		}
		// Notify.
		err = c.Notifier.NotifyGroupUpdated(group)
		if err != nil {
			return meh.Wrap(err, "notify", meh.Details{"group": group})
		}
		return nil
	})
	if err != nil {
		return meh.Wrap(err, "run in tx", nil)
	}
	return nil
}

// DeleteGroupByID deletes the group with the given id and notifies.
func (c *Controller) DeleteGroupByID(ctx context.Context, groupID uuid.UUID) error {
	err := pgutil.RunInTx(ctx, c.DB, func(ctx context.Context, tx pgx.Tx) error {
		// Delete in store.
		err := c.Store.DeleteGroupByID(ctx, tx, groupID)
		if err != nil {
			return meh.Wrap(err, "delete group in store", meh.Details{"group_id": groupID})
		}
		// Notify.
		err = c.Notifier.NotifyGroupDeleted(groupID)
		if err != nil {
			return meh.Wrap(err, "notify", meh.Details{"group_id": groupID})
		}
		return nil
	})
	if err != nil {
		return meh.Wrap(err, "run in tx", nil)
	}
	return nil
}

// GroupByID retrieves a store.Group by its id.
func (c *Controller) GroupByID(ctx context.Context, groupID uuid.UUID) (store.Group, error) {
	var group store.Group
	err := pgutil.RunInTx(ctx, c.DB, func(ctx context.Context, tx pgx.Tx) error {
		var err error
		group, err = c.Store.GroupByID(ctx, tx, groupID)
		if err != nil {
			return meh.Wrap(err, "group from store", meh.Details{"group_id": groupID})
		}
		return nil
	})
	if err != nil {
		return store.Group{}, meh.Wrap(err, "run in tx", nil)
	}
	return group, nil
}

// Groups retrieves a paginated store.Group list.
func (c *Controller) Groups(ctx context.Context, filters store.GroupFilters, params pagination.Params) (pagination.Paginated[store.Group], error) {
	var groups pagination.Paginated[store.Group]
	err := pgutil.RunInTx(ctx, c.DB, func(ctx context.Context, tx pgx.Tx) error {
		var err error
		groups, err = c.Store.Groups(ctx, tx, filters, params)
		if err != nil {
			return meh.Wrap(err, "groups from store", meh.Details{
				"filters":           filters,
				"pagination_params": params,
			})
		}
		return nil
	})
	if err != nil {
		return pagination.Paginated[store.Group]{}, meh.Wrap(err, "run in tx", nil)
	}
	return groups, nil
}
