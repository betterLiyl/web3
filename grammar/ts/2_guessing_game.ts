import * as readline from 'readline';

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

function getRandomNumber(min: number, max: number): number {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

function askQuestion(question: string): Promise<string> {
    return new Promise((resolve) => {
        rl.question(question, (answer) => {
            resolve(answer);
        });
    });
}

async function guessingGame(): Promise<void> {
    console.log('欢迎来到猜数字游戏！');
    console.log('我已经想好了一个1到100之间的数字，请猜猜看！');
    
    const secretNumber = getRandomNumber(1, 100);
    
    while (true) {
        const input = await askQuestion('请输入你的猜测: ');
        const guess = parseInt(input.trim());
        
        if (isNaN(guess)) {
            console.log('请输入一个有效的数字！');
            continue;
        }
        
        if (guess < secretNumber) {
            console.log('太小了！再试试看。');
        } else if (guess > secretNumber) {
            console.log('太大了！再试试看。');
        } else {
            console.log(`恭喜你！你猜对了！答案就是 ${secretNumber}`);
            break;
        }
    }
    
    rl.close();
}

guessingGame().catch(console.error);