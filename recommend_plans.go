/*
	コードサンプル
	空いている予定を5つ提案する
*/

package main

import (
	"fmt"
	"os"
	"time"
)

const (
	RCD_NUM     = 5                           //レコメンド予定数
	TIME_FORMAT = "2006-01-02 15:04:05 -0700" //goでは日付フォーマットは特定の日時を指定する仕様
)

//予定の構造体
type Plan struct {
	start, end string //開始日時と終了日時
}

//開始時刻をUnixTimeで返す
func (p *Plan) GetStartUnixTime() int64 {
	return StringToUnixtime(p.start)
}

//終了時刻をUnixTimeで返す
func (p *Plan) GetEndUnixTime() int64 {
	return StringToUnixtime(p.end)
}

//日付文字列をtime型に変換し返却
func StringToUnixtime(str_time string) int64 {
	t, err := time.Parse(TIME_FORMAT, str_time+" +0900")
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	return t.Unix()
}

/*
  空き予定を返却する ※商用では期間制約を設ける
  plans 予定のあるPlan構造体配列(予定開始時間でソートされていること)
  recommend_start_date 空きを探す開始時間(YYYY-MM-DD HH:MM:SS)
  interval_min 確保したい時間(分)

  確保したい時間分の終了時間が予定と被っていないか確認していく。
  被っていない場合は提案時間とし提案数に達するまで確認し続ける。
  指定された開始時刻から早い時間を提案していくので時間を進めながら確認していく。

  処理概要は、9:00-9:30が空いているか確認したい場合、9:30が予定にないかをチェックしていく。
  9:00-10:00で予定がある場合は、10:30に時間を進め10:00-10:30で予定が空いているかチェックする。
  終了時間と開始時間の重複は許容され、10:00まで予定としても10:00-10:30を提案する
*/
func GetRecommendPlans(plans []Plan, recommend_start_date string, interval_min int) []Plan {
	recommend_plans := make([]Plan, RCD_NUM) //提案予定格納用配列

	var interval_sec int64 = int64(interval_min * 60)                            //UnixTimeで処理するので確保時間(分)を秒に変換
	var check_time int64 = StringToUnixtime(recommend_start_date) + interval_sec //空き予定確認 - 終了時間
	var pos_rcd_plan int = 0                                                     //レコメンド予定配列の位置
	var pos_plan int = 0                                                         //予定配列の位置

	for pos_rcd_plan < RCD_NUM { //レコメンド数に達したら終了
		fmt.Println("-------")
		fmt.Println("空き確認時間(の終了時間):", time.Unix(check_time, 0))

		if pos_plan < len(plans) {
			//既存の予定と被っていたら空き確認時間を進める
			if check_time > plans[pos_plan].GetStartUnixTime() && check_time <= plans[pos_plan].GetEndUnixTime() {
				fmt.Println("この時間は予定あり=>", plans[pos_plan])
				check_time = plans[pos_plan].GetEndUnixTime() + interval_sec //次の空き確認時間として既存予定の終了時間にする
				pos_plan++                                                   //次の予定に進める
				continue
			}
		}
		//レコメンド予定設定
		recommend_plans[pos_rcd_plan].start = time.Unix(check_time-interval_sec, 0).String()
		recommend_plans[pos_rcd_plan].end = time.Unix(check_time, 0).String()
		check_time = check_time + interval_sec //次の空き確認時間として今回空いていた時間の終了時間にする
		fmt.Println("この時間帯は空きあり！", recommend_plans[pos_rcd_plan].start, recommend_plans[pos_rcd_plan].end)
		pos_rcd_plan++
	}
	return recommend_plans
}

/*
  確定している予定を返却
  ※今回は固定で返す
  ※商用では期間制約を設ける
  ※予定開始順に並んでる配列を返す
*/
func GetPlans() []Plan {
	var plans = []Plan{
		{start: "2015-12-19 09:00:00", end: "2015-12-19 10:00:00"},
		{"2015-12-19 11:00:00", "2015-12-19 12:00:00"},
		{"2015-12-19 12:00:00", "2015-12-19 12:30:00"},
		{"2015-12-19 13:00:00", "2015-12-19 15:00:00"},
		{"2015-12-19 15:00:00", "2015-12-19 19:00:00"},
	}
	return plans
}

func main() {
	//予定取得 範囲指定はされた方がいい
	plans := GetPlans()

	//レコメンド時間
	recommend_plans := GetRecommendPlans(plans, "2015-12-19 09:00:00", 30)
	//レコメンド予定表示
	fmt.Println("*************************************")
	for i := 0; i < len(recommend_plans); i++ {
		fmt.Println(recommend_plans[i])
	}
}
