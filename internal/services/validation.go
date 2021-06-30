package services

import (
	"context"
	"fmt"

	"github.com/authzed/spicedb/internal/namespace"
	v0 "github.com/authzed/spicedb/pkg/proto/authzed/api/v0"
)

type invalidRelationError struct {
	error
	subject *v0.User
	onr     *v0.ObjectAndRelation
}

func validateTupleWrite(ctx context.Context, tpl *v0.RelationTuple, nsm namespace.Manager) error {
	if err := nsm.CheckNamespaceAndRelation(
		ctx,
		tpl.ObjectAndRelation.Namespace,
		tpl.ObjectAndRelation.Relation,
		false, // Disallow ellipsis
	); err != nil {
		return err
	}

	if err := nsm.CheckNamespaceAndRelation(
		ctx,
		tpl.User.GetUserset().Namespace,
		tpl.User.GetUserset().Relation,
		true, // Allow Ellipsis
	); err != nil {
		return err
	}

	_, ts, _, err := nsm.ReadNamespaceAndTypes(ctx, tpl.ObjectAndRelation.Namespace)
	if err != nil {
		return err
	}

	isAllowed, err := ts.IsAllowedDirectRelation(
		tpl.ObjectAndRelation.Relation,
		tpl.User.GetUserset().Namespace,
		tpl.User.GetUserset().Relation)
	if err != nil {
		return err
	}

	if isAllowed == namespace.DirectRelationNotValid {
		return invalidRelationError{
			error:   fmt.Errorf("relation %v is not allowed on the right hand side of %v", tpl.User, tpl.ObjectAndRelation),
			subject: tpl.User,
			onr:     tpl.ObjectAndRelation,
		}
	}

	return nil
}