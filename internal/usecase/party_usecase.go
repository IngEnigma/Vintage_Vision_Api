package usecase

import (
	"context"
	"errors"
	"time"
	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/utils"
)

type PartyUsecase struct {
	Repo domain.PartyRepository
}

func NewPartyUsecase(repo domain.PartyRepository) *PartyUsecase {
	return &PartyUsecase{Repo: repo}
}

func (p *PartyUsecase) CreateParty(ctx context.Context, party *domain.Party) error {
	party.PartyCode = utils.GenerateUniqueCode()
	party.ExpiresAt = time.Now().Add(3 * time.Hour)

	err := p.Repo.Create(ctx, party)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgCreateParty, err)
		return err
	}

	utils.Logger.Infof("%s: %s", constants.MsgPartyCreatedSuccessfully, party.PartyCode)
	return nil
}

func (p *PartyUsecase) JoinParty(ctx context.Context, code string, profileID uint) error {
	party, err := p.Repo.FindByCode(ctx, code)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgPartyNotFound, err)
		return err
	}

	isMember, err := p.Repo.IsMember(ctx, party.ID, profileID)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgCheckPartyMembership, err)
		return err
	}
	if isMember {
		utils.Logger.Warnf("ProfileID %d already in party %d", profileID, party.ID)
		return errors.New(constants.ErrMsgAlreadyInParty)
	}

	err = p.Repo.AddMember(ctx, party.ID, profileID)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgAddPartyMember, err)
		return err
	}

	utils.Logger.Infof("Joined: PartyID %d, ProfileID %d", party.ID, profileID)
	return nil
}

func (p *PartyUsecase) GetPartyByCode(ctx context.Context, code string) (*domain.Party, error) {
	party, err := p.Repo.FindByCode(ctx, code)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgFindPartyByCode, err)
		return nil, err
	}
	utils.Logger.Infof("Retrieved: %s", party.PartyCode)
	return party, nil
}

func (p *PartyUsecase) LeaveParty(ctx context.Context, partyID, profileID uint) error {
	err := p.Repo.RemoveMember(ctx, partyID, profileID)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgRemovePartyMember, err)
		return err
	}
	utils.Logger.Infof("Left: PartyID %d, ProfileID %d", partyID, profileID)
	return nil
}
