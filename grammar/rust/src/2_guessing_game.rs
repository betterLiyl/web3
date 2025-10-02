use rand::Rng;
use std::cmp::Ordering;
use std::io;

fn main() {
    println!("欢迎来到猜数字游戏！");
    println!("我已经想好了一个1到100之间的数字，请猜猜看！");

    let secret_number = rand::thread_rng().gen_range(1..=100);

    loop {
        println!("请输入你的猜测:");

        let mut guess = String::new();

        io::stdin()
            .read_line(&mut guess)
            .expect("读取输入失败");

        let guess: u32 = match guess.trim().parse() {
            Ok(num) => num,
            Err(_) => {
                println!("请输入一个有效的数字！");
                continue;
            }
        };

        match guess.cmp(&secret_number) {
            Ordering::Less => println!("太小了！再试试看。"),
            Ordering::Greater => println!("太大了！再试试看。"),
            Ordering::Equal => {
                println!("恭喜你！你猜对了！答案就是 {}", secret_number);
                break;
            }
        }
    }
}