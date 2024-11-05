package entity

import (
	pb "github.com/nakamaFramework/cgp-common/proto"
	"sort"
)

type Hand struct {
	userId string
	cards  []*pb.Card
}

func NewHand(userId string, cards []*pb.Card) *Hand {
	return &Hand{
		userId: userId,
		cards:  cards,
	}
}

func NewHandFromPb(v *pb.ShanPlayerHand) *Hand {
	return &Hand{
		userId: v.UserId,
		cards:  v.Hand.Cards,
	}
}

func (h *Hand) ToPb() *pb.ShanPlayerHand {
	point, handType := h.Eval()
	return &pb.ShanPlayerHand{
		UserId: h.userId,
		Hand: &pb.ShanHand{
			Cards: h.cards,
			Point: point,
			Type:  handType,
		},
	}
}

func getCardPoint(r pb.CardRank) int32 {
	switch r {
	case pb.CardRank_RANK_A:
		return 1
	case pb.CardRank_RANK_2:
		return 2
	case pb.CardRank_RANK_3:
		return 3
	case pb.CardRank_RANK_4:
		return 4
	case pb.CardRank_RANK_5:
		return 5
	case pb.CardRank_RANK_6:
		return 6
	case pb.CardRank_RANK_7:
		return 7
	case pb.CardRank_RANK_8:
		return 8
	case pb.CardRank_RANK_9:
		return 9
	case pb.CardRank_RANK_10, pb.CardRank_RANK_J, pb.CardRank_RANK_Q, pb.CardRank_RANK_K:
		return 0
	default:
		return 0
	}
}

func calculatePoint(cards []*pb.Card) int32 {
	if cards == nil {
		return 0
	}
	point := int32(0)
	for _, c := range cards {
		point += getCardPoint(c.Rank)
	}
	return point
}

func (h *Hand) Eval() (int32, pb.ShanHandType) {
	point := calculatePoint(h.cards)

	// Kiểm tra bài "Shan" (2 lá đầu tiên)
	if len(h.cards) == 2 && (point == 8 || point == 9) {
		return point, pb.ShanHandType_SHAN
	}

	// Kiểm tra bài "Xám" (3 lá giống nhau)
	if len(h.cards) == 3 && isXam(h.cards) {
		return point, pb.ShanHandType_XAM
	}

	// Kiểm tra bài "3 con đầu người" (J, Q, K)
	if len(h.cards) == 3 && isThreeFace(h.cards) {
		return point, pb.ShanHandType_THREE_FACE
	}

	// Kiểm tra "Thùng Phá Sảnh"
	if len(h.cards) == 3 && isFlushStraight(h.cards) {
		return point, pb.ShanHandType_FLUSH_STRAIGHT
	}

	// Nếu không có gì đặc biệt, trả về điểm và loại bình thường
	return point, pb.ShanHandType_NORMAL
}

func isSameSuit(cards []*pb.Card) bool {
	if len(cards) < 2 {
		return false
	}
	suit := cards[0].Suit
	for _, card := range cards {
		if card.Suit != suit {
			return false
		}
	}
	return true
}

func isPair(cards []*pb.Card) bool {
	if len(cards) == 2 {
		return cards[0].Rank == cards[1].Rank
	}
	return false
}

func isXam(cards []*pb.Card) bool {
	if len(cards) != 3 {
		return false
	}
	rank := cards[0].Rank
	for _, card := range cards {
		if card.Rank != rank {
			return false
		}
	}
	return true
}

func isThreeFace(cards []*pb.Card) bool {
	if len(cards) != 3 {
		return false
	}
	for _, card := range cards {
		if card.Rank != pb.CardRank_RANK_J && card.Rank != pb.CardRank_RANK_Q && card.Rank != pb.CardRank_RANK_K {
			return false
		}
	}
	return true
}
func isFlushStraight(cards []*pb.Card) bool {
	// Kiểm tra số lượng lá bài và cùng chất
	if len(cards) != 3 {
		return false
	}

	// Kiểm tra xem tất cả các lá bài có cùng chất không
	suit := cards[0].Suit
	for _, card := range cards[1:] {
		if card.Suit != suit {
			return false
		}
	}

	// Chuyển các lá bài thành slice các giá trị rank
	ranks := []int{int(cards[0].Rank), int(cards[1].Rank), int(cards[2].Rank)}
	sort.Ints(ranks)

	// Kiểm tra nếu các lá bài là A-2-3 hoặc K-A-2
	if isSpecialCase(ranks) {
		return false
	}

	// Kiểm tra tính liên tiếp của các lá bài
	return ranks[2]-ranks[0] == 2
}

// Kiểm tra các trường hợp đặc biệt A-2-3 hoặc K-A-2
func isSpecialCase(ranks []int) bool {
	return (ranks[0] == int(pb.CardRank_RANK_A) && ranks[1] == int(pb.CardRank_RANK_2) && ranks[2] == int(pb.CardRank_RANK_3)) ||
		(ranks[0] == int(pb.CardRank_RANK_K) && ranks[1] == int(pb.CardRank_RANK_A) && ranks[2] == int(pb.CardRank_RANK_2))
}

func (h *Hand) PlayerCanDraw() bool {
	// Tính điểm của 2 lá đầu tiên
	if len(h.cards) == 2 && calculatePoint(h.cards) < 8 {
		return true // Cho phép rút lá thứ 3 nếu điểm dưới 8
	}
	return false // Không cho phép rút nếu đã đạt 8 hoặc 9
}

func (h *Hand) AddCards(c []*pb.Card) {
	h.cards = append(h.cards, c...)
}

// comparing player hand with dealer hand, -1 -> lost, 1 -> win, 0 -> tie
func (h *Hand) Compare(d *Hand) int {
	playerPoint, playerType := h.Eval()
	dealerPoint, dealerType := d.Eval()

	// So sánh loại tay bài trước
	if playerType != dealerType {
		if playerType > dealerType {
			return 1 // Người chơi thắng
		}
		return -1 // Nhà cái thắng
	}

	// Nếu loại tay bài giống nhau, so sánh điểm
	if playerPoint != dealerPoint {
		if playerPoint > dealerPoint {
			return 1
		}
		return -1
	}

	// Nếu điểm bằng nhau, so sánh số lá bài (tay bài có ít lá hơn thắng)
	if len(h.cards) != len(d.cards) {
		if len(h.cards) < len(d.cards) {
			return 1
		}
		return -1
	}

	// Nếu cùng số lá bài, so sánh lá bài cao nhất
	maxPlayerCard := getMaxCard(h.cards)
	maxDealerCard := getMaxCard(d.cards)
	if maxPlayerCard.Rank != maxDealerCard.Rank {
		if maxPlayerCard.Rank > maxDealerCard.Rank {
			return 1
		}
		return -1
	}
	// Nếu cùng số lá và lá bài cao nhất có cùng Rank, so sánh chất
	if maxPlayerCard.Suit > maxDealerCard.Suit {
		return 1
	} else if maxPlayerCard.Suit < maxDealerCard.Suit {
		return -1
	}
	return 0 // Hòa
}

func getMaxCard(cards []*pb.Card) *pb.Card {
	if len(cards) == 0 {
		return nil
	}
	maxCard := cards[0]
	for _, card := range cards[1:] {
		if card.Rank > maxCard.Rank || (card.Rank == maxCard.Rank && card.Suit > maxCard.Suit) {
			maxCard = card
		}
	}
	return maxCard
}
