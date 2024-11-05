// package main

// import (
// 	"fmt"
// 	"minh-shan-plus-module/usecase/engine"
// )

// func main() {
// 	gameEngine := engine.NewGameEngine()
// 	for {
// 		gameEngine.StartGame()

// 		// Hỏi người chơi có muốn tiếp tục không
// 		var input string
// 		fmt.Println("Bạn có muốn chơi ván mới không? (y/n)")
// 		fmt.Scanln(&input)

//			if input == "y" {
//				gameEngine.ResetGame()
//			} else {
//				fmt.Println("Kết thúc trò chơi. Cảm ơn bạn đã chơi!")
//				break
//			}
//		}
//	}
package main

import (
	"context"
	"database/sql"
	"time"
	"minh-shan-plus-module/api"
	"minh-shan-plus-module/entity"
	"github.com/heroiclabs/nakama-common/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

// protojson.MarshalOptions và protojson.UnmarshalOptions 
// được sử dụng để chuyển đổi dữ liệu giữa định dạng JSON và protobuf, 
//điều này gợi ý rằng đầu vào/đầu ra có thể liên quan đến định dạng protobuf.

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	initStart := time.Now()
	marshaler := &protojson.MarshalOptions{
		UseEnumNumbers:  true,
		EmitUnpopulated: true,
	}
	unmarshaler := &protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}
	if err := initializer.RegisterMatch(entity.ModuleName, func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		return api.NewMatchHandler(marshaler, unmarshaler), nil
	}); err != nil {
		return err
	}
	logger.Info("Plugin loaded in '%d' msec.", time.Since(initStart).Milliseconds())
	return nil
	
}
