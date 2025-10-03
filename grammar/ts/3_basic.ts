// 全局变量和常量
let globalVar: number = 100;
const globalString: string = "全局字符串";
const PI: number = 3.14159;
const MAX_SIZE: number = 1000;

// 类型别名
type ID = string | number;
type Status = "pending" | "completed" | "failed";

// 接口定义
interface Person {
    name: string;
    age: number;
    city: string;
    introduce(): string;
}

// 扩展接口
interface Employee extends Person {
    employeeId: ID;
    department: string;
    salary?: number; // 可选属性
}

// 泛型接口
interface Container<T> {
    value: T;
    getValue(): T;
    setValue(value: T): void;
}

// 类定义
class PersonClass implements Person {
    name: string;
    age: number;
    city: string;
    
    // 构造函数
    constructor(name: string, age: number, city: string) {
        this.name = name;
        this.age = age;
        this.city = city;
    }
    
    // 方法
    introduce(): string {
        return `我叫${this.name}，今年${this.age}岁，来自${this.city}`;
    }
    
    // 设置年龄
    setAge(age: number): void {
        this.age = age;
    }
    
    // 静态方法
    static createDefault(): PersonClass {
        return new PersonClass("默认姓名", 0, "默认城市");
    }
}

// 抽象类
abstract class Animal {
    protected name: string;
    
    constructor(name: string) {
        this.name = name;
    }
    
    // 抽象方法
    abstract makeSound(): string;
    
    // 具体方法
    getName(): string {
        return this.name;
    }
}

// 继承抽象类
class Dog extends Animal {
    private breed: string;
    
    constructor(name: string, breed: string) {
        super(name);
        this.breed = breed;
    }
    
    makeSound(): string {
        return "汪汪！";
    }
    
    getBreed(): string {
        return this.breed;
    }
}

// 枚举
enum Color {
    Red = "red",
    Green = "green",
    Blue = "blue"
}

enum Direction {
    Up,
    Down,
    Left,
    Right
}

// 联合类型
type Theme = "light" | "dark";
type Size = "small" | "medium" | "large";

// 普通函数
function add(a: number, b: number): number {
    return a + b;
}

// 箭头函数
const multiply = (a: number, b: number): number => a * b;

// 可选参数和默认参数
// 必需参数: 函数定义时直接声明的参数
// 可选参数: 在参数名后添加问号 ?
// 默认参数: 在参数赋值时指定默认值
function greet(name: string, greeting: string = "你好", punctuation?: string): string {
    return `${greeting}, ${name}${punctuation || "!"}`;
}

// 剩余参数
function sum(...numbers: number[]): number {
    return numbers.reduce((total, num) => total + num, 0);
}

// 函数重载
function processValue(value: string): string;
function processValue(value: number): number;
function processValue(value: boolean): boolean;
function processValue(value: string | number | boolean): string | number | boolean {
    if (typeof value === "string") {
        return value.toUpperCase();
    } else if (typeof value === "number") {
        return value * 2;
    } else {
        return !value;
    }
}

// 泛型函数
function identity<T>(arg: T): T {
    return arg;
}

function findMax<T>(array: T[], compareFn: (a: T, b: T) => number): T | undefined {
    if (array.length === 0) return undefined;
    
    let max = array[0];
    for (let i = 1; i < array.length; i++) {
        if (compareFn(array[i], max) > 0) {
            max = array[i];
        }
    }
    return max;
}

// 高阶函数
function createMultiplier(factor: number): (x: number) => number {
    return (x: number) => x * factor;
}

// 演示变量声明
function variableDemo(): void {
    console.log("=== 变量声明演示 ===");
    
    // let 和 const
    let mutableVar: number = 10;
    const immutableVar: number = 20;
    
    console.log(`可变变量: ${mutableVar}, 不可变变量: ${immutableVar}`);
    
    mutableVar = 30; // 可以修改
    // immutableVar = 40; // 编译错误
    
    // 类型推断
    let inferredString = "这是推断的字符串类型";
    let inferredNumber = 42;
    let inferredBoolean = true;
    
    console.log(`推断类型: ${typeof inferredString}, ${typeof inferredNumber}, ${typeof inferredBoolean}`);
    
    // 数组
    let numbers: number[] = [1, 2, 3, 4, 5];
    let strings: Array<string> = ["a", "b", "c"];
    
    console.log(`数字数组: ${numbers}, 字符串数组: ${strings}`);
    
    // 元组
    let tuple: [string, number, boolean] = ["hello", 42, true];
    console.log(`元组: ${tuple[0]}, ${tuple[1]}, ${tuple[2]}`);
    
    // 对象
    let person: { name: string; age: number } = {
        name: "张三",
        age: 25
    };
    console.log(`对象: ${person.name}, ${person.age}岁`);
    
    // 联合类型
    let unionVar: string | number = "hello";
    console.log(`联合类型初始值: ${unionVar}`);
    unionVar = 42;
    console.log(`联合类型修改后: ${unionVar}`);
    
    // 字面量类型
    let theme: Theme = "light";
    let size: Size = "medium";
    console.log(`主题: ${theme}, 尺寸: ${size}`);
    
    // null 和 undefined
    let nullableVar: string | null = null;
    let undefinedVar: string | undefined = undefined;
    console.log(`可空变量: ${nullableVar}, 未定义变量: ${undefinedVar}`);
    
    // 全局变量使用
    console.log(`全局变量: ${globalVar}, 全局字符串: ${globalString}`);
    console.log(`常量: PI=${PI}, MAX_SIZE=${MAX_SIZE}`);
}

// 演示函数
function functionDemo(): void {
    console.log("\n=== 函数演示 ===");
    
    // 调用普通函数
    const result = add(5, 3);
    console.log(`5 + 3 = ${result}`);
    
    // 调用箭头函数
    const product = multiply(4, 6);
    console.log(`4 * 6 = ${product}`);
    
    // 调用带默认参数的函数
    console.log(greet("李四"));
    console.log(greet("王五", "早上好"));
    console.log(greet("赵六", "晚上好", "。"));
    
    // 调用剩余参数函数
    const total = sum(1, 2, 3, 4, 5);
    console.log(`1+2+3+4+5 = ${total}`);
    
    // 调用重载函数
    console.log(`字符串处理: ${processValue("hello")}`);
    console.log(`数字处理: ${processValue(21)}`);
    console.log(`布尔处理: ${processValue(false)}`);
    
    // 调用泛型函数
    console.log(`泛型函数: ${identity("TypeScript")}`);
    console.log(`泛型函数: ${identity(42)}`);
    
    const maxNumber = findMax([3, 7, 1, 9, 2], (a, b) => a - b);
    console.log(`最大数字: ${maxNumber}`);
    
    // 高阶函数
    const double = createMultiplier(2);
    const triple = createMultiplier(3);
    console.log(`双倍函数: ${double(5)}, 三倍函数: ${triple(5)}`);
    
    // 匿名函数和立即执行函数
    const anonymousResult = ((x: number, y: number) => x + y)(10, 20);
    console.log(`匿名函数结果: ${anonymousResult}`);
}

// 演示类和接口
function classInterfaceDemo(): void {
    console.log("\n=== 类和接口演示 ===");
    
    // 创建类实例
    const person1 = new PersonClass("李四", 28, "北京");
    const person2 = PersonClass.createDefault();
    
    console.log(person1.introduce());
    console.log(person2.introduce());
    
    // 修改属性
    person1.setAge(29);
    console.log(`${person1.name} 更新年龄后: ${person1.age}岁`);
    
    // 继承示例
    const dog = new Dog("旺财", "金毛");
    console.log(`动物名字: ${dog.getName()}`);
    console.log(`狗的叫声: ${dog.makeSound()}`);
    console.log(`狗的品种: ${dog.getBreed()}`);
    
    // 接口使用
    const employee: Employee = {
        name: "张三",
        age: 30,
        city: "上海",
        employeeId: "EMP001",
        department: "技术部",
        salary: 15000,
        introduce() {
            return `我是${this.name}，员工编号${this.employeeId}，在${this.department}工作`;
        }
    };
    
    console.log(employee.introduce());
    
    // 泛型类使用
    class GenericContainer<T> implements Container<T> {
        value: T;
        
        constructor(value: T) {
            this.value = value;
        }
        
        getValue(): T {
            return this.value;
        }
        
        setValue(value: T): void {
            this.value = value;
        }
    }
    
    const stringContainer = new GenericContainer<string>("Hello TypeScript");
    const numberContainer = new GenericContainer<number>(42);
    
    console.log(`字符串容器: ${stringContainer.getValue()}`);
    console.log(`数字容器: ${numberContainer.getValue()}`);
}

// 演示循环
function loopDemo(): void {
    console.log("\n=== 循环演示 ===");
    
    // for 循环
    console.log("传统for循环:");
    for (let i = 1; i <= 5; i++) {
        console.log(`  第${i}次循环`);
    }
    
    // for...of 循环 - 数组
    const numbers = [10, 20, 30, 40, 50];
    console.log("for...of遍历数组:");
    numbers.forEach((value, index) => {
        console.log(`  [${index}] = ${value}`);
    });
    
    // for...of 循环 - 字符串
    const text = "Hello";
    console.log("for...of遍历字符串:");
    for (let i = 0; i < text.length; i++) {
        console.log(`  ${i}: ${text[i]}`);
    }
    
    // for...in 循环 - 对象
    const person = { name: "张三", age: 25, city: "北京" };
    console.log("for...in遍历对象:");
    for (const key in person) {
        console.log(`  ${key}: ${person[key as keyof typeof person]}`);
    }
    
    // while 循环
    console.log("while循环:");
    let count = 1;
    while (count <= 3) {
        console.log(`  while: ${count}`);
        count++;
    }
    
    // do...while 循环
    console.log("do...while循环:");
    let num = 1;
    do {
        console.log(`  do-while: ${num}`);
        num++;
    } while (num <= 3);
    
    // forEach 方法
    console.log("forEach方法:");
    numbers.forEach((value, index) => {
        console.log(`  forEach[${index}] = ${value}`);
    });
    
    // map, filter, reduce 示例
    const doubled = numbers.map(x => x * 2);
    const evens = numbers.filter(x => x % 20 === 0);
    const sum = numbers.reduce((acc, curr) => acc + curr, 0);
    
    console.log(`原数组: ${numbers}`);
    console.log(`双倍: ${doubled}`);
    console.log(`能被20整除的: ${evens}`);
    console.log(`总和: ${sum}`);
}

// 演示条件判断
function conditionDemo(): void {
    console.log("\n=== 条件判断演示 ===");
    
    const age = 18;
    const score = 85;
    
    // if-else 基本用法
    if (age >= 18) {
        console.log(`年龄${age}岁，已成年`);
    } else {
        console.log(`年龄${age}岁，未成年`);
    }
    
    // if-else if-else
    if (score >= 90) {
        console.log(`分数${score}，等级：优秀`);
    } else if (score >= 80) {
        console.log(`分数${score}，等级：良好`);
    } else if (score >= 60) {
        console.log(`分数${score}，等级：及格`);
    } else {
        console.log(`分数${score}，等级：不及格`);
    }
    
    // 三元运算符
    const status = age >= 18 ? "成年人" : "未成年人";
    console.log(`状态: ${status}`);
    
    // 逻辑运算符
    const temperature = 25;
    const humidity = 60;
    if (temperature > 20 && humidity < 70) {
        console.log(`温度${temperature}°C，湿度${humidity}%，天气舒适`);
    }
    
    // 短路求值
    const name = "";
    const displayName = name || "匿名用户";
    console.log(`显示名称: ${displayName}`);
    
    // 空值合并运算符
    const userInput: string | null = null;
    const defaultValue = userInput ?? "默认值";
    console.log(`空值合并结果: ${defaultValue}`);
    
    // 可选链
    const user: { profile?: { name?: string } } = {};
    const userName = user.profile?.name ?? "未知用户";
    console.log(`用户名: ${userName}`);
    
    // 类型守卫
    function isString(value: unknown): value is string {
        return typeof value === "string";
    }
    
    const unknownValue: unknown = "Hello TypeScript";
    if (isString(unknownValue)) {
        console.log(`类型守卫确认为字符串: ${unknownValue.toUpperCase()}`);
    }
}

// 演示switch语句
function switchDemo(): void {
    console.log("\n=== Switch演示 ===");
    
    // 基本switch
    const day: number = 3;
    let dayName: string;
    switch (day) {
        case 1:
            dayName = "星期一";
            break;
        case 2:
            dayName = "星期二";
            break;
        case 3:
            dayName = "星期三";
            break;
        case 4:
            dayName = "星期四";
            break;
        case 5:
            dayName = "星期五";
            break;
        case 6:
        case 7:
            dayName = "周末";
            break;
        default:
            dayName = "无效的日期";
    }
    console.log(`第${day}天是: ${dayName}`);
    
    // switch 表达式（使用函数）
    function getTimeOfDay(hour: number): string {
        switch (true) {
            case hour >= 0 && hour < 6:
                return "深夜";
            case hour >= 6 && hour < 12:
                return "上午";
            case hour >= 12 && hour < 18:
                return "下午";
            case hour >= 18 && hour < 24:
                return "晚上";
            default:
                return "无效时间";
        }
    }
    
    const hour = 14;
    console.log(`${hour}点属于: ${getTimeOfDay(hour)}`);
    
    // switch 枚举 - 使用函数参数避免类型推断问题
    function handleColor(color: Color): void {
        switch (color) {
            case Color.Red:
                console.log("选择了红色");
                break;
            case Color.Green:
                console.log("选择了绿色");
                break;
            case Color.Blue:
                console.log("选择了蓝色");
                break;
            default:
                console.log("未知颜色");
        }
    }
    
    handleColor(Color.Red);
    handleColor(Color.Green);
    handleColor(Color.Blue);
    
    // switch 联合类型 - 使用函数参数避免类型推断问题
    function handleTheme(theme: Theme): void {
        switch (theme) {
            case "light":
                console.log("浅色主题");
                break;
            case "dark":
                console.log("深色主题");
                break;
            default:
                // TypeScript 会检查是否处理了所有情况
                const exhaustiveCheck: never = theme;
                throw new Error(`未处理的主题: ${exhaustiveCheck}`);
        }
    }
    
    handleTheme("dark");
    handleTheme("light");
    
    // switch 类型判断
    function handleValue(value: string | number | boolean): string {
        switch (typeof value) {
            case "string":
                return `字符串: ${value.toUpperCase()}`;
            case "number":
                return `数字: ${value * 2}`;
            case "boolean":
                return `布尔值: ${!value}`;
            default:
                return "未知类型";
        }
    }
    
    console.log(handleValue("hello"));
    console.log(handleValue(42));
    console.log(handleValue(true));
}

// 演示异步编程
async function asyncDemo(): Promise<void> {
    console.log("\n=== 异步编程演示 ===");
    
    // Promise
    const promise = new Promise<string>((resolve, reject) => {
        setTimeout(() => {
            resolve("Promise 完成");
        }, 1000);
    });
    
    promise.then(result => {
        console.log(`Promise 结果: ${result}`);
    });
    
    // async/await
    async function fetchData(): Promise<string> {
        return new Promise(resolve => {
            setTimeout(() => resolve("异步数据"), 500);
        });
    }
    
    try {
        const data = await fetchData();
        console.log(`Async/await 结果: ${data}`);
    } catch (error) {
        console.error("异步错误:", error);
    }
    
    // 并行异步操作
    const [result1, result2] = await Promise.all([
        fetchData(),
        new Promise<number>(resolve => setTimeout(() => resolve(42), 300))
    ]);
    
    console.log(`并行结果: ${result1}, ${result2}`);
}

// 主函数
async function main(): Promise<void> {
    console.log("TypeScript语言基础语法学习");
    console.log("==========================");
    
    // 演示各个语法特性
    variableDemo();
    functionDemo();
    classInterfaceDemo();
    loopDemo();
    conditionDemo();
    switchDemo();
    //await 关键字用于等待异步操作完成
    // await asyncDemo();
    asyncDemo();
    
    console.log("\n学习完成！");
}

// 运行主函数
main().catch(console.error);