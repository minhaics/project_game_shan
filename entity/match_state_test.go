package entity

import (
	"fmt"
	"testing"

	pb "github.com/nakamaFramework/cgp-common/proto"
)


func TestMatchState(t *testing.T) {
	// Tạo một MatchState mới để kiểm tra các thao tác
	s := NewMatchState(&MatchLabel{
		Open:     5,
		Bet:      MaxBetAllowed,
		Code:     "test",
		Name:     "test_table",
		Password: "",
		MaxSize:  5,
	})
	s.Init()
	s.PlayingPresences.Put("A", FakePrecense{})//
	s.PlayingPresences.Put("B", FakePrecense{})

	// Tạo và xáo bài
	deck := NewDeck()
	deck.Shuffle()

	// Thêm cược của người chơi "A" với số tiền 100
	s.AddBet(&pb.BlackjackBet{
		UserId: "A",
		Chips:  100,
	})
	// Thêm cược của người chơi "A" với số tiền 200
	s.AddBet(&pb.BlackjackBet{
		UserId: "A",
		Chips:  200,
	})

	// Kiểm tra xem người chơi "A" đã đặt cược hay chưa
	// Nếu `IsBet("A")` trả về `false`, tức là có lỗi trong logic đặt cược
	if !s.IsBet("A") {
		t.Errorf("Expected player A to have placed a bet, but IsBet returned false")
	}

	// Kiểm tra tổng tiền cược của người chơi "A"
	// Ở đây, kỳ vọng rằng tổng cược là 300 (100 + 200), nếu khác, sẽ có lỗi
	totalBet := s.GetTotalBet("A") // Giả sử hàm này đã được định nghĩa
	if totalBet != 300 {
		t.Errorf("Expected total bet for player A to be 300, but got %d", totalBet)
	}

	// Chia 2 lá bài cho ngân hàng (banker) và thêm chúng vào trạng thái ván bài
	if cards, err := deck.Deal(2); err != nil {
		t.Fatalf(err.Error())
	} else {
		s.AddCards(cards.Cards, "", pb.BlackjackHandN0_BLACKJACK_HAND_UNSPECIFIED)
	}

	// Chia 2 lá bài cho người chơi "A" và thêm chúng vào tay của họ
	if cards, err := deck.Deal(2); err != nil {
		t.Fatalf(err.Error())
	} else {
		s.AddCards(cards.Cards, "A", pb.BlackjackHandN0_BLACKJACK_HAND_1ST)
	}

	// Tính toán kết quả cuối cùng của ván chơi và kiểm tra kết quả
	// Giả sử `CalcGameFinish` trả về kết quả mà ta có thể so sánh
	result := s.CalcGameFinish()
	expectedResult := /* Giá trị mong đợi */ nil // Thay `nil` bằng giá trị bạn mong đợi từ `CalcGameFinish`
	if result != expectedResult {
		t.Errorf("Expected game result to be %v, but got %v", expectedResult, result)
	}

	// In ra kết quả cuối cùng (có thể không cần thiết trong bài test)
	fmt.Printf("====END GAME====\n%v\n", result)
}
