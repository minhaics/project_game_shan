package engine

import (
	"minh-shan-plus-module/entity"

	pb "github.com/nakamaFramework/cgp-common/proto"
	"google.golang.org/protobuf/proto"
)

type UseCase interface {
	NewGame(s *entity.MatchState) error
	Deal(amount int) []*pb.Card
	Finish(s *entity.MatchState) *pb.BlackjackUpdateFinish
	Draw(s *entity.MatchState, userId string, handN0 pb.BlackjackHandN0)
	DoubleDown(s *entity.MatchState, userId string, handN0 pb.BlackjackHandN0) int64
	Split(s *entity.MatchState, userId string) int64
	Insurance(s *entity.MatchState, userId string) int64
	RejoinUserMessage(s *entity.MatchState, userId string) map[pb.OpCodeUpdate]proto.Message
}
