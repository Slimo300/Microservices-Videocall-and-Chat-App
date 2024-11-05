package orm

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	merrors "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *GroupsGormRepository) GetMemberByID(ctx context.Context, userID, memberID uuid.UUID) (models.Member, error) {
	var m Member
	if err := r.db.WithContext(ctx).First(&m, memberID).Error; err != nil {
		return models.Member{}, merrors.NewMemberNotFoundError(memberID.String())
	}
	var issuer Member
	if err := r.db.WithContext(ctx).Where(Member{UserID: userID, GroupID: m.GroupID}).First(&issuer).Error; err != nil {
		return models.Member{}, merrors.NewUserNotInGroupError(userID.String(), m.GroupID.String())
	}
	member := unmarshalMember(m)
	return member, nil
}

func (r *GroupsGormRepository) DeleteMember(ctx context.Context, userID, memberID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var i, t Member
		if err := tx.First(&t, memberID).Error; err != nil {
			return merrors.NewMemberNotFoundError(memberID.String())
		}
		if err := tx.Where(Member{UserID: userID, GroupID: t.GroupID}).First(&i).Error; err != nil {
			return merrors.NewUserNotInGroupError(userID.String(), t.GroupID.String())
		}
		issuer := unmarshalMember(i)
		target := unmarshalMember(t)
		if !issuer.CanDelete(target) {
			return merrors.NewMemberUnauthorizedError(target.GroupID().String(), merrors.DeleteMemberAction())
		}
		return tx.Delete(&t).Error
	})
}

func (r *GroupsGormRepository) UpdateMember(ctx context.Context, userID, memberID uuid.UUID, updateFn func(m *models.Member) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var i, t Member
		if err := tx.First(&t, memberID).Error; err != nil {
			return merrors.NewMemberNotFoundError(memberID.String())
		}
		if err := tx.Where(Member{UserID: userID, GroupID: t.GroupID}).First(&i).Error; err != nil {
			return merrors.NewUserNotInGroupError(userID.String(), t.GroupID.String())
		}
		issuer := unmarshalMember(i)
		target := unmarshalMember(t)
		if !issuer.CanAlter(target) {
			return merrors.NewMemberUnauthorizedError(t.GroupID.String(), merrors.UpdateMemberAction())
		}
		if err := updateFn(&target); err != nil {
			return err
		}
		t = marshalMember(target)
		return tx.Save(&t).Error
	})
}
