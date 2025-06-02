package repository

import (
	"context"
	"errors"
	"time"
	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/utils"
	"vintage-vision-api/repository/entity"
	"vintage-vision-api/repository/mapper"

	"gorm.io/gorm"
)

type PartyRepo struct {
	DB *gorm.DB
}

func NewPartyRepo(db *gorm.DB) domain.PartyRepository {
	return &PartyRepo{DB: db}
}

func (r *PartyRepo) Create(ctx context.Context, party *domain.Party) error {
	entityParty := mapper.FromDomainParty(party)
	err := r.DB.WithContext(ctx).Create(&entityParty).Error
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgCreateParty, err)
		return err
	}

	party.ID = entityParty.ID
	utils.Logger.Infof("%s: ID %d", constants.MsgPartyCreatedSuccessfully, party.ID)
	return nil
}

func (r *PartyRepo) FindByCode(ctx context.Context, code string) (*domain.Party, error) {
	var entityParty entity.Party
	err := r.DB.WithContext(ctx).Where("code = ?", code).First(&entityParty).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Logger.Warnf("%s: %s", constants.ErrMsgPartyNotFound, code)
			return nil, err
		}
		utils.Logger.Errorf("%s: %v", constants.ErrMsgFindPartyByCode, err)
		return nil, err
	}

	domainParty := mapper.ToDomainParty(&entityParty)
	utils.Logger.Infof("%s: %s", constants.MsgPartyRetrievedSuccessfully, domainParty.PartyCode)
	return domainParty, nil
}

func (r *PartyRepo) IsMember(ctx context.Context, partyID, profileID uint) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&entity.PartyMember{}).
		Where("party_id = ? AND profile_id = ?", partyID, profileID).
		Count(&count).Error
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgCheckPartyMembership, err)
		return false, err
	}

	isMember := count > 0
	utils.Logger.Infof("ID %d: %t", partyID, isMember)
	return isMember, nil
}

func (r *PartyRepo) AddMember(ctx context.Context, partyID, profileID uint) error {
	if isMember, err := r.IsMember(ctx, partyID, profileID); err != nil {
		return err
	} else if isMember {
		utils.Logger.Warnf("ID %d is already in party ID %d", profileID, partyID)
		return nil
	}

	member := mapper.FromDomainPartyMember(&domain.PartyMember{
		PartyID:   partyID,
		ProfileID: profileID,
		JoinedAt:  time.Now(),
	})

	err := r.DB.WithContext(ctx).Create(&member).Error
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgAddPartyMember, err)
		return err
	}

	utils.Logger.Infof("ID %d added to party ID %d", profileID, partyID)
	return nil
}

func (r *PartyRepo) RemoveMember(ctx context.Context, partyID, profileID uint) error {
	if isMember, err := r.IsMember(ctx, partyID, profileID); err != nil {
		return err
	} else if !isMember {
		utils.Logger.Warnf("ID %d is not a member of party ID %d", profileID, partyID)
		return nil
	}

	err := r.DB.WithContext(ctx).Where("party_id = ? AND profile_id = ?", partyID, profileID).Delete(&entity.PartyMember{}).Error
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgRemovePartyMember, err)
		return err
	}

	utils.Logger.Infof("ID %d removed from ID %d", profileID, partyID)
	return nil
}
