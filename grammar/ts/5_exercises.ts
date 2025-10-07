// TypeScript 语言练习题
// 请完成以下练习题，每个函数都有详细的说明和示例
// 重点练习 TypeScript 的类型系统、接口、泛型、类等核心概念

// ==================== 练习题 ====================

// 练习1: 基础类型和类型注解
// 完成函数，计算两个数字的和、差、积、商
interface CalculationResult {
    sum: number;
    difference: number;
    product: number;
    quotient: number;
}

function calculate(a: number, b: number): CalculationResult {
    // TODO: 实现计算逻辑
    // 返回包含四种运算结果的对象
    return {
        sum: a + b,
        difference: a - b,
        product: a * b,
        quotient: a / b
    }; 
}

// 练习2: 字符串操作和联合类型
type StringOperation = "uppercase" | "lowercase" | "capitalize" | "reverse";

function processString(input: string, operation: StringOperation): string {
    // TODO: 根据操作类型处理字符串
    switch (operation) {
        case "uppercase":
            return input.toUpperCase();
        case "lowercase":
            return input.toLowerCase();
        case "capitalize":
            return input.charAt(0).toUpperCase() + input.slice(1);
        case "reverse":
            return input.split("").reverse().join("");
        default:
            throw new Error("Invalid operation");
    }
}

// 练习3: 数组操作和泛型
// 完成泛型函数，找出数组中的最大值和最小值
interface MinMax<T> {
    min: T;
    max: T;
}

function findMinMax<T>(array: T[], compareFn: (a: T, b: T) => number): MinMax<T> | null {
    // TODO: 实现找出最大值和最小值的逻辑
    // 如果数组为空，返回 null
    // 使用 compareFn 进行比较
    if (array.length === 0) {
        return null;
    }
    let min = array[0];
    let max = array[0];
    for (let i = 1; i < array.length; i++) {
        if (compareFn(array[i], min) < 0) {
            min = array[i];
        }
        if (compareFn(array[i], max) > 0) {
            max = array[i];
        }
    }
    return { min, max };
}

// 练习4: 接口和类
// 定义一个学生管理系统
interface Student {
    // TODO: 定义学生接口
    id: number;
    name: string;
    age: number;
    grades: number[];
    getAverageGrade(): number;
    addGrade(grade: number): void;
}

class StudentClass implements Student {
    // TODO: 实现 Student 接口
    id: number = 0;
    name: string = "";
    age: number = 0;
    grades: number[] = [];

    constructor(id: number, name: string, age: number) {
        // TODO: 初始化属性
        this.id = id;
        this.name = name;
        this.age = age;
    }

    getAverageGrade(): number {
        // TODO: 计算平均成绩
        if (this.grades.length === 0) {
            return 0;
        }
        const sum = this.grades.reduce((acc, cur) => acc + cur, 0);
        return sum / this.grades.length;
    }

    addGrade(grade: number): void {
        // TODO: 添加成绩到数组
        this.grades.push(grade);
    }

    // TODO: 添加一个方法获取学生信息
    getInfo(): string {
        return `ID: ${this.id}, Name: ${this.name}, Age: ${this.age}, Grades: ${this.grades}`;
    }
}

// 练习5: 枚举和类型守卫
enum TaskStatus {
    // TODO: 定义任务状态枚举
    PENDING = "pending",
    IN_PROGRESS = "in_progress",
    COMPLETED = "completed",
    CANCELLED = "cancelled"
}

interface Task {
    id: number;
    title: string;
    status: TaskStatus;
    createdAt: Date;
    completedAt?: Date;
}

// 类型守卫函数
function isCompleted(task: Task): task is Task & { completedAt: Date } {
    // TODO: 实现类型守卫，检查任务是否已完成
    return task.status === TaskStatus.COMPLETED;
}

function updateTaskStatus(task: Task, newStatus: TaskStatus): Task {
    // TODO: 更新任务状态
    // 如果状态变为 COMPLETED，设置 completedAt
    // 如果状态从 COMPLETED 变为其他状态，清除 completedAt
    if (newStatus !== TaskStatus.COMPLETED && isCompleted(task)) {
        task.completedAt = new Date();
    }
    task.status = newStatus;
    return task;
}

// 练习6: 泛型约束和映射类型
// 定义一个泛型函数，只接受有 length 属性的类型
interface Lengthwise {
    length: number;
}

function logLength<T extends Lengthwise>(arg: T): T {
    // TODO: 打印参数的长度并返回参数
    console.log(`Length: ${arg.length}`);
    return arg;
}

// 映射类型练习：创建一个类型，将所有属性变为可选
type Partial<T> = {
    // TODO: 实现 Partial 映射类型
    // [P in keyof T]?: T[P];
    [P in keyof T]?: T[P];
};

// 练习7: 函数重载
// 实现一个格式化函数，支持不同类型的输入
function format(value: string): string;
function format(value: number): string;
function format(value: Date): string;
function format(value: boolean): string;
function format(value: string | number | Date | boolean): string {
    // TODO: 根据不同类型实现格式化逻辑
    // string: 添加引号
    // number: 添加千位分隔符
    // Date: 格式化为 YYYY-MM-DD
    // boolean: 转换为 Yes/No
    if (typeof value === "string") {
        return `"${value}"`;
    } else if (typeof value === "number") {
        return value.toLocaleString();
    } else if (value instanceof Date) {
        return value.toISOString().split("T")[0];
    } else if (typeof value === "boolean") {
        return value ? "Yes" : "No";
    }
    return "";
}

// 练习8: 装饰器（实验性功能）
// 创建一个简单的日志装饰器
function log(target: any, propertyName: string, descriptor: PropertyDescriptor) {
    // TODO: 实现日志装饰器
    // 在方法执行前后打印日志
    const method = descriptor.value;
    
    descriptor.value = function (...args: any[]) {
        // TODO: 添加日志逻辑
        console.log(`调用方法: ${propertyName}`);
        const result = method.apply(this, args);
        console.log(`方法 ${propertyName} 执行完成`);
        return result;
    };
}

class Calculator {
    @log
    add(a: number, b: number): number {
        return a + b;
    }

    @log
    multiply(a: number, b: number): number {
        return a * b;
    }
}

// 练习9: Promise 和异步操作
// 实现一个异步数据获取函数
interface ApiResponse<T> {
    data: T;
    status: number;
    message: string;
}

async function fetchData<T>(url: string): Promise<ApiResponse<T>> {
    // TODO: 模拟异步数据获取
    // 随机延迟 1-3 秒
    // 90% 概率成功，10% 概率失败
    await new Promise(resolve => setTimeout(resolve, Math.random() * 2000 + 1000));
    // 90% 概率成功，10% 概率失败
    if (Math.random() < 0.9) {
        return {
            data: {} as T,
            status: 200,
            message: "Success"
        };
    } else {
        throw new Error("模拟 API 调用失败");
    }
}

// 练习10: 高级类型操作
// 实现一个深度只读类型
type DeepReadonly<T> = {
    // TODO: 实现深度只读类型
    // readonly [P in keyof T]: T[P] extends object ? DeepReadonly<T[P]> : T[P];
    readonly [P in keyof T]: T[P] extends object ? DeepReadonly<T[P]> : T[P];
};

// 实现一个选择特定属性的类型
type Pick<T, K extends keyof T> = {
    // TODO: 实现 Pick 类型
    // [P in K]: T[P];
    [P in K]: T[P];
};

// ==================== 测试函数 ====================

function testCalculate(): void {
    console.log("=== 测试计算函数 ===");
    const result = calculate(10, 3);
    console.log("10 和 3 的计算结果:", result);
    // 期望结果: { sum: 13, difference: 7, product: 30, quotient: 3.33... }
}

function testStringProcessing(): void {
    console.log("\n=== 测试字符串处理 ===");
    const text = "hello world";
    console.log(`原文: ${text}`);
    console.log(`大写: ${processString(text, "uppercase")}`);
    console.log(`小写: ${processString(text, "lowercase")}`);
    console.log(`首字母大写: ${processString(text, "capitalize")}`);
    console.log(`反转: ${processString(text, "reverse")}`);
}

function testMinMax(): void {
    console.log("\n=== 测试最大最小值 ===");
    const numbers = [3, 7, 2, 9, 1, 5, 8];
    const result = findMinMax(numbers, (a, b) => a - b);
    console.log(`数组 [${numbers.join(", ")}] 的最大最小值:`, result);
    
    const strings = ["apple", "banana", "cherry", "date"];
    const stringResult = findMinMax(strings, (a, b) => a.localeCompare(b));
    console.log(`字符串数组的最大最小值:`, stringResult);
}

function testStudent(): void {
    console.log("\n=== 测试学生类 ===");
    // TODO: 创建学生实例并测试方法
    const student = new StudentClass(1, "张三", 20);
    student.addGrade(85);
    student.addGrade(92);
    student.addGrade(78);
    console.log(student.getInfo());
    console.log(`平均成绩: ${student.getAverageGrade()}`);
}

function testTaskManagement(): void {
    console.log("\n=== 测试任务管理 ===");
    // TODO: 创建任务并测试状态更新
    const task: Task = {
        id: 1,
        title: "完成 TypeScript 练习",
        status: TaskStatus.PENDING,
        createdAt: new Date()
    };
    
    console.log("初始任务:", task);
    const updatedTask = updateTaskStatus(task, TaskStatus.COMPLETED);
    console.log("更新后任务:", updatedTask);
    console.log("任务是否完成:", isCompleted(updatedTask));
}

function testGenericConstraints(): void {
    console.log("\n=== 测试泛型约束 ===");
    logLength("Hello World");
    logLength([1, 2, 3, 4, 5]);
    logLength({ length: 10, value: "test" });
}

function testFunctionOverloads(): void {
    console.log("\n=== 测试函数重载 ===");
    console.log(`字符串格式化: ${format("hello")}`);
    console.log(`数字格式化: ${format(1234567)}`);
    console.log(`日期格式化: ${format(new Date())}`);
    console.log(`布尔值格式化: ${format(true)}`);
}

function testDecorator(): void {
    console.log("\n=== 测试装饰器 ===");
    const calc = new Calculator();
    console.log("加法结果:", calc.add(5, 3));
    console.log("乘法结果:", calc.multiply(4, 6));
}

async function testAsyncOperations(): Promise<void> {
    console.log("\n=== 测试异步操作 ===");
    try {
        const response = await fetchData<{ name: string; age: number }>("/api/user");
        console.log("API 响应:", response);
    } catch (error) {
        console.error("API 调用失败:", error);
    }
}

// 主函数
async function main(): Promise<void> {
    console.log("TypeScript 语言练习题");
    console.log("请完成上面标有 TODO 的函数实现");
    console.log("然后运行测试函数查看结果");
    console.log("=====================================");
    
    // 运行测试函数
    testCalculate();
    testStringProcessing();
    testMinMax();
    testStudent();
    testTaskManagement();
    testGenericConstraints();
    testFunctionOverloads();
    testDecorator();
    await testAsyncOperations();
    
    console.log("\n练习完成后，你可以尝试以下进阶挑战：");
    console.log("1. 实现一个类型安全的状态管理器");
    console.log("2. 创建一个支持插件系统的应用框架");
    console.log("3. 使用装饰器实现依赖注入系统");
    console.log("4. 实现一个类型安全的 ORM 查询构建器");
    console.log("5. 创建一个支持类型推断的配置系统");
}

// 运行主函数
main().catch(console.error);

// 导出类型和接口供其他模块使用
export {
    CalculationResult,
    StringOperation,
    MinMax,
    Student,
    StudentClass,
    TaskStatus,
    Task,
    ApiResponse,
    DeepReadonly,
    Pick
};