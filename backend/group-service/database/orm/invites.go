package orm

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	merrors "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *GroupsGormRepository) GetInviteByID(ctx context.Context, userID, inviteID uuid.UUID) (models.Invite, error) {
	var i Invite
	if err := r.db.WithContext(ctx).Preload("Iss").Preload("Group").Preload("Target").First(&i, inviteID).Error; err != nil {
		return models.Invite{}, err
	}
	invite := unmarshalInvite(i)
	if !invite.CanUserSee(userID) {
		return models.Invite{}, merrors.NewUserCantSeeInviteError(userID.String(), inviteID.String())
	}
	return unmarshalInvite(i), nil
}

func (r *GroupsGormRepository) GetUserInvites(ctx context.Context, userID uuid.UUID, num, offset int) ([]models.Invite, error) {
	var is []Invite
	if err := r.db.WithContext(ctx).Order("modified DESC").Limit(num).Offset(offset).
		Where(Invite{TargetID: userID}).
		Or(Invite{IssId: userID}).
		Preload("Iss").Preload("Group").Preload("Target").Find(&is).Error; err != nil {
		return nil, err
	}
	return unmarshalInvites(is), nil
}

func (r *GroupsGormRepository) CreateInvite(ctx context.Context, invite models.Invite) (models.Invite, error) {
	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var m Member
		if err := tx.Where(Member{GroupID: invite.GroupID(), UserID: invite.IssuerID()}).First(&m).Error; err != nil {
			return merrors.NewUserNotInGroupError(invite.IssuerID().String(), invite.GroupID().String())
		}
		member := unmarshalMember(m)
		if !member.CanSendInvite() {
			return merrors.NewMemberUnauthorizedError(invite.GroupID().String(), merrors.AddMemberAction())
		}
		if err := tx.First(&User{}, invite.TargetID()).Error; err != nil {
			return merrors.NewUserNotFoundError(invite.TargetID().String())
		}
		if err := tx.Where(Member{GroupID: invite.GroupID(), UserID: invite.TargetID()}).First(&Member{}).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
			return merrors.NewUserAlreadyInGroupError(invite.TargetID().String(), invite.GroupID().String())
		}
		if err := tx.Where(Invite{TargetID: invite.TargetID(), GroupID: invite.GroupID(), Status: int(models.INVITE_AWAITING)}).First(&Invite{}).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
			return merrors.NewUserAlreadyInvitedError(invite.TargetID().String(), invite.GroupID().String())
		}
		i := marshalInvite(invite)
		if err := tx.Create(&i).Error; err != nil {
			return err
		}
		if err := tx.Preload("Target").Preload("Iss").Preload("Group").First(&i, invite.ID()).Error; err != nil {
			return err
		}
		invite = unmarshalInvite(i)
		return nil
	}); err != nil {
		return models.Invite{}, err
	}
	return invite, nil
}

func (r *GroupsGormRepository) UpdateInvite(ctx context.Context, inviteID uuid.UUID, updateFn func(i *models.Invite) (*models.Member, error)) error {
	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var i Invite
		if err := tx.First(&i, inviteID).Error; err != nil {
			return merrors.NewInviteNotFoundError(inviteID.String())
		}
		invite := unmarshalInvite(i)
		member, err := updateFn(&invite)
		if err != nil {
			return err
		}
		i = marshalInvite(invite)
		if err := tx.Save(&i).Error; err != nil {
			return err
		}
		if member == nil {
			return nil
		}
		m := marshalMember(*member)
		return tx.Create(m).Error
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}); err != nil {
		return err
	}
	return nil
}
