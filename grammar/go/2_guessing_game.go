package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main2() {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())
	
	// 生成1-100之间的随机数
	secretNumber := rand.Intn(100) + 1
	fmt.Println("欢迎来到猜数字游戏！")
	fmt.Println("我已经想好了一个1到100之间的数字，请猜猜看！")
	
	for {
		fmt.Print("请输入你的猜测: ")
		var guess int
		_, err := fmt.Scanf("%d", &guess)
		
		if err != nil {
			fmt.Println("请输入一个有效的数字！")
			continue
		}
		
		if guess < secretNumber {
			fmt.Println("太小了！再试试看。")
		} else if guess > secretNumber {
			fmt.Println("太大了！再试试看。")
		} else {
			fmt.Printf("恭喜你！你猜对了！答案就是 %d\n", secretNumber)
			break
		}
	}
}