module CardRank
  RANK_A = 1
  RANK_2 = 2
  RANK_3 = 3
  RANK_4 = 4
  RANK_5 = 5
  RANK_6 = 6
  RANK_7 = 7
  RANK_8 = 8
  RANK_9 = 9
  RANK_10 = 10
  RANK_J = 11
  RANK_Q = 12
  RANK_K = 13
end

module ShanHandType
  NORMAL = 0
  SHAN = 1
  XAM = 2
  THREE_FACE = 3
  FLUSH_STRAIGHT = 4
end

class Card
  attr_accessor :suit, :rank

  def initialize(suit, rank)
    @suit = suit
    @rank = rank
  end
end

class ShanPlayerHand
    attr_accessor :user_id, :cards
    def initialize(user_id, cards)
        @user_id = user_id
        @cards = cards
    end

  def get_card_point(rank)
    case rank
    when CardRank::RANK_A
      1
    when CardRank::RANK_2..CardRank::RANK_9
      rank
    else
      0
    end
  end

  def calculate_point
    @cards.sum { |card| get_card_point(card.rank) }
  end

  def eval
    point = calculate_point

    if @cards.size == 2 && (point == 8 || point == 9)
      return [point, ShanHandType::SHAN]
    end
    if @cards.size == 3
      return [point, ShanHandType::XAM] if is_xam
      return [point, ShanHandType::THREE_FACE] if is_three_face
      return [point, ShanHandType::FLUSH_STRAIGHT] if is_flush_straight
    end

    [point, ShanHandType::NORMAL]
  end

  def is_same_suit
    @cards.all? { |card| card.suit == @cards[0].suit }
  end

  def is_xam
    @cards.map(&:rank).uniq.size == 1
  end

  def is_three_face
    @cards.all? { |card| [CardRank::RANK_J, CardRank::RANK_Q, CardRank::RANK_K].include?(card.rank) }
  end

  def is_flush_straight
    return false unless @cards.size == 3 && is_same_suit

    ranks = @cards.map(&:rank).sort
    return false if ranks == [CardRank::RANK_A, CardRank::RANK_2, CardRank::RANK_3] || ranks == [CardRank::RANK_A, CardRank::RANK_2, CardRank::RANK_K]
    
    ranks[2] - ranks[0] == 2
  end
end

# Ví dụ sử dụng
cards = [
  Card.new("hearts", CardRank::RANK_A),
  Card.new("hearts", CardRank::RANK_2),
  Card.new("hearts", CardRank::RANK_3)
]

hand = ShanPlayerHand.new("player1", cards)
point, hand_type = hand.eval
puts "Điểm: #{point}, Loại bài: #{hand_type}"
