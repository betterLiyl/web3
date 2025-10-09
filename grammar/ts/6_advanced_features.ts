// ============= 1. 泛型 (Generics) =============

console.log("=== 泛型基础 ===");

// 基本泛型函数
function identity<T>(arg: T): T {
    return arg;
}

let output1 = identity<string>("myString");
let output2 = identity<number>(100);
let output3 = identity("myString"); // 类型推断

console.log("泛型函数结果:", output1, output2, output3);

// 泛型接口
interface GenericIdentityFn<T> {
    (arg: T): T;
}

let myIdentity: GenericIdentityFn<number> = identity;

// 泛型类
class GenericNumber<T> {
    zeroValue: T;
    add: (x: T, y: T) => T;

    constructor(zeroValue: T, addFn: (x: T, y: T) => T) {
        this.zeroValue = zeroValue;
        this.add = addFn;
    }
}

let myGenericNumber = new GenericNumber<number>(0, (x, y) => x + y);
console.log("泛型类结果:", myGenericNumber.add(myGenericNumber.zeroValue, 5));

let stringNumeric = new GenericNumber<string>("", (x, y) => x + y);
console.log("字符串泛型:", stringNumeric.add(stringNumeric.zeroValue, "test"));

// 泛型约束
interface Lengthwise {
    length: number;
}

function loggingIdentity<T extends Lengthwise>(arg: T): T {
    console.log("长度:", arg.length);
    return arg;
}

loggingIdentity("hello");
loggingIdentity([1, 2, 3]);
loggingIdentity({ length: 10, value: 3 });

// 在泛型约束中使用类型参数
function getProperty<T, K extends keyof T>(obj: T, key: K): T[K] {
    return obj[key];
}

let obj1 = { a: 1, b: 2, c: 3, d: 4 };
console.log("获取属性:", getProperty(obj1, "a"));

// ============= 2. 高级类型 =============

console.log("\n=== 高级类型 ===");

// 联合类型
type StringOrNumber = string | number;

function padLeft(value: string, padding: StringOrNumber) {
    if (typeof padding === "number") {
        return Array(padding + 1).join(" ") + value;
    }
    if (typeof padding === "string") {
        return padding + value;
    }
    throw new Error(`Expected string or number, got '${padding}'.`);
}

console.log("联合类型:", padLeft("Hello world", 4));

// 交叉类型
interface ErrorHandling {
    success: boolean;
    error?: { message: string };
}

interface ArtworksData {
    artworks: { title: string }[];
}

interface ArtistsData {
    artists: { name: string }[];
}

type ArtworksResponse = ArtworksData & ErrorHandling;
type ArtistsResponse = ArtistsData & ErrorHandling;

const handleArtistsResponse = (response: ArtistsResponse) => {
    if (response.error) {
        console.error(response.error.message);
        return;
    }
    console.log("艺术家:", response.artists);
};

// 类型守卫
function isFish(pet: Fish | Bird): pet is Fish {
    return (pet as Fish).swim !== undefined;
}

interface Fish {
    swim(): void;
}

interface Bird {
    fly(): void;
}

function move(pet: Fish | Bird) {
    if (isFish(pet)) {
        pet.swim();
    } else {
        pet.fly();
    }
}

// 可辨识联合
interface Square {
    kind: "square";
    size: number;
}

interface Rectangle {
    kind: "rectangle";
    width: number;
    height: number;
}

interface Circle {
    kind: "circle";
    radius: number;
}

type Shape = Square | Rectangle | Circle;

function area(s: Shape): number {
    switch (s.kind) {
        case "square": return s.size * s.size;
        case "rectangle": return s.height * s.width;
        case "circle": return Math.PI * s.radius ** 2;
    }
}

console.log("正方形面积:", area({ kind: "square", size: 5 }));

// 索引类型
function pluck<T, K extends keyof T>(o: T, propertyNames: K[]): T[K][] {
    return propertyNames.map(n => o[n]);
}

interface Car {
    manufacturer: string;
    model: string;
    year: number;
}

let taxi: Car = {
    manufacturer: 'Toyota',
    model: 'Camry',
    year: 2014
};

let makeAndModel: string[] = pluck(taxi, ['manufacturer', 'model']);
console.log("索引类型结果:", makeAndModel);

// 映射类型
type MyReadonly<T> = {
    readonly [P in keyof T]: T[P];
};

type MyPartial<T> = {
    [P in keyof T]?: T[P];
};

interface UserData {
    id: number;
    name: string;
    email: string;
}

type ReadonlyUser = MyReadonly<UserData>;
type PartialUser = MyPartial<UserData>;

const readonlyUser: ReadonlyUser = { id: 1, name: "Alice", email: "alice@example.com" };
const partialUser: PartialUser = { name: "Bob" };

console.log("映射类型示例:", readonlyUser, partialUser);

type CarPartial = MyPartial<Car>;
type CarReadonly = MyReadonly<Car>;

// 条件类型
type TypeName<T> =
    T extends string ? "string" :
    T extends number ? "number" :
    T extends boolean ? "boolean" :
    T extends undefined ? "undefined" :
    T extends Function ? "function" :
    "object";

type TypeName0 = TypeName<string>;  // "string"
type TypeName1 = TypeName<"a">;     // "string"
type TypeName2 = TypeName<true>;    // "boolean"
type TypeName3 = TypeName<() => void>; // "function"
type TypeName4 = TypeName<string[]>; // "object"

// ============= 3. 装饰器 (Decorators) =============

console.log("\n=== 装饰器 ===");

// 类装饰器
function sealed(constructor: Function) {
    Object.seal(constructor);
    Object.seal(constructor.prototype);
}

function classDecorator<T extends { new(...args: any[]): {} }>(constructor: T) {
    return class extends constructor {
        newProperty = "new property";
        hello = "override";
    }
}

@classDecorator
@sealed
class Greeter {
    greeting: string;
    constructor(message: string) {
        this.greeting = message;
    }
    greet() {
        return "Hello, " + this.greeting;
    }
}

console.log("装饰器类:", new Greeter("world"));

// 方法装饰器
function enumerable(value: boolean) {
    return function (target: any, propertyKey: string, descriptor: PropertyDescriptor) {
        descriptor.enumerable = value;
    };
}

function log(target: any, propertyName: string, descriptor: PropertyDescriptor) {
    const method = descriptor.value;
    descriptor.value = function (...args: any[]) {
        console.log(`调用方法 ${propertyName}，参数:`, args);
        const result = method.apply(this, args);
        console.log(`方法 ${propertyName} 返回:`, result);
        return result;
    }
}

class Calculator {
    @log
    @enumerable(false)
    add(a: number, b: number): number {
        return a + b;
    }
}

const calc = new Calculator();
calc.add(2, 3);

// 属性装饰器
function format(formatString: string) {
    return function (target: any, propertyKey: string) {
        let value = target[propertyKey];

        const getter = function () {
            return `${formatString} ${value}`;
        };

        const setter = function (newVal: string) {
            value = newVal;
        };

        Object.defineProperty(target, propertyKey, {
            get: getter,
            set: setter,
            enumerable: true,
            configurable: true
        });
    };
}

class Person {
    @format("Hello")
    name: string;

    constructor(name: string) {
        this.name = name;
    }
}

const person = new Person("Alice");
console.log("属性装饰器:", person.name);

// 参数装饰器
function required(target: any, propertyKey: string, parameterIndex: number) {
    console.log(`参数装饰器 - 方法: ${propertyKey}, 参数索引: ${parameterIndex}`);
}

class User {
    greet(@required name: string): string {
        return `Hello ${name}`;
    }
}

// ============= 4. 模块系统 =============

console.log("\n=== 模块系统 ===");

// 导出
export interface StringValidator {
    isAcceptable(s: string): boolean;
}

export const numberRegexp = /^[0-9]+$/;

export class ZipCodeValidator implements StringValidator {
    isAcceptable(s: string) {
        return s.length === 5 && numberRegexp.test(s);
    }
}

// 重新导出
export { ZipCodeValidator as mainValidator };

// 默认导出
export default class DefaultValidator implements StringValidator {
    isAcceptable(s: string) {
        return s.length > 0;
    }
}

// 命名空间
namespace Validation {
    export interface StringValidator {
        isAcceptable(s: string): boolean;
    }

    const lettersRegexp = /^[A-Za-z]+$/;
    const numberRegexp = /^[0-9]+$/;

    export class LettersOnlyValidator implements StringValidator {
        isAcceptable(s: string) {
            return lettersRegexp.test(s);
        }
    }

    export class ZipCodeValidator implements StringValidator {
        isAcceptable(s: string) {
            return s.length === 5 && numberRegexp.test(s);
        }
    }
}

// 使用命名空间
let validators: { [s: string]: Validation.StringValidator; } = {};
validators["ZIP code"] = new Validation.ZipCodeValidator();
validators["Letters only"] = new Validation.LettersOnlyValidator();

// ============= 5. 异步编程 =============

console.log("\n=== 异步编程 ===");

// Promise
function delay(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function asyncExample() {
    console.log("开始异步操作");
    await delay(1000);
    console.log("异步操作完成");
}

// 泛型 Promise
function fetchData<T>(url: string): Promise<T> {
    return fetch(url).then(response => response.json() as T);
}

// async/await 与错误处理
async function fetchUserData(id: number): Promise<{ name: string; email: string } | null> {
    try {
        const response = await fetch(`/api/users/${id}`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return await response.json() as { name: string; email: string };
    } catch (error) {
        console.error("获取用户数据失败:", error);
        return null;
    }
}

// Promise 组合
async function parallelRequests() {
    const [user1, user2, user3] = await Promise.all([
        fetchUserData(1),
        fetchUserData(2),
        fetchUserData(3)
    ]);
    
    console.log("并行请求结果:", { user1, user2, user3 });
}

// ============= 6. 工具类型 =============

console.log("\n=== 工具类型 ===");

interface Todo {
    title: string;
    description: string;
    completed: boolean;
}

// Partial<T>
function updateTodo(todo: Todo, fieldsToUpdate: Partial<Todo>) {
    return { ...todo, ...fieldsToUpdate };
}

const todo1: Todo = {
    title: "学习 TypeScript",
    description: "学习 TypeScript 高级特性",
    completed: false,
};

const todo2 = updateTodo(todo1, {
    completed: true,
});

console.log("Partial 示例:", todo2);

// Required<T>
type RequiredTodo = Required<Todo>;

// Readonly<T>
const readonlyTodo: Readonly<Todo> = {
    title: "只读任务",
    description: "这个任务是只读的",
    completed: false,
};

// Pick<T, K>
type TodoPreview = Pick<Todo, "title" | "completed">;

const preview: TodoPreview = {
    title: "预览任务",
    completed: false,
};

// Omit<T, K>
type TodoInfo = Omit<Todo, "completed">;

const todoInfo: TodoInfo = {
    title: "信息任务",
    description: "只包含标题和描述",
};

// Record<K, T>
type Page = "home" | "about" | "contact";

const nav: Record<Page, { title: string; url: string }> = {
    home: { title: "首页", url: "/" },
    about: { title: "关于", url: "/about" },
    contact: { title: "联系", url: "/contact" },
};

console.log("Record 示例:", nav);

// Exclude<T, U>
type ExcludeExample = Exclude<"a" | "b" | "c", "a">;  // "b" | "c"

// Extract<T, U>
type ExtractExample = Extract<"a" | "b" | "c", "a" | "f">;  // "a"

// NonNullable<T>
type NonNullableExample = NonNullable<string | number | undefined>;  // string | number

// ReturnType<T>
function f1(): { a: number; b: string } {
    return { a: 1, b: "hello" };
}

type ReturnTypeExample = ReturnType<typeof f1>;  // { a: number; b: string }

// ============= 7. 类型操作 =============

console.log("\n=== 类型操作 ===");

// keyof 操作符
type Point = { x: number; y: number };
type P = keyof Point;  // "x" | "y"

function getPropertyFromPoint<T, K extends keyof T>(obj: T, key: K): T[K] {
    return obj[key];
}

const point: Point = { x: 1, y: 2 };
const pointX = getPropertyFromPoint(point, "x");  // number
console.log("keyof 示例:", pointX);

// typeof 操作符
let s = "hello";
let n: typeof s;  // string

const person2 = { name: "Alice", age: 30 };
type PersonType = typeof person2;  // { name: string; age: number }

function createPerson(): PersonType {
    return { name: "Bob", age: 25 };
}

const newPerson = createPerson();
console.log("typeof 示例:", newPerson);

// 索引访问类型
type Age = PersonType["age"];  // number
type NameOrAge = PersonType["name" | "age"];  // string | number
type PersonKeys = PersonType[keyof PersonType];  // string | number

// 模板字面量类型
type World = "world";
type Greeting = `hello ${World}`;  // "hello world"

type EmailLocaleIDs = "welcome_email" | "email_heading";
type FooterLocaleIDs = "footer_title" | "footer_sendoff";

type AllLocaleIDs = `${EmailLocaleIDs | FooterLocaleIDs}_id`;
// "welcome_email_id" | "email_heading_id" | "footer_title_id" | "footer_sendoff_id"

// ============= 8. 声明合并 =============

console.log("\n=== 声明合并 ===");

// 接口合并
interface Box {
    height: number;
    width: number;
}

interface Box {
    scale: number;
}

let box: Box = { height: 5, width: 6, scale: 10 };
console.log("接口合并:", box);

// 命名空间合并
namespace Animals {
    export class Zebra { }
}

namespace Animals {
    export interface Legged {
        numberOfLegs: number;
    }
    export class Dog { }
}

// 等同于:
// namespace Animals {
//     export interface Legged {
//         numberOfLegs: number;
//     }
//     export class Zebra { }
//     export class Dog { }
// }

// ============= 9. 三斜线指令 =============

/// <reference path="..." />
/// <reference types="..." />
/// <reference lib="..." />

// ============= 10. 实用示例 =============

console.log("\n=== 实用示例 ===");

// 深度只读
type DeepReadonly<T> = {
    readonly [P in keyof T]: T[P] extends object ? DeepReadonly<T[P]> : T[P];
};

interface NestedObject {
    a: {
        b: {
            c: string;
        };
    };
}

type ReadonlyNested = DeepReadonly<NestedObject>;

// 函数重载 (需要DOM类型支持)
// function createElement(tagName: "img"): HTMLImageElement;
// function createElement(tagName: "input"): HTMLInputElement;
// function createElement(tagName: string): Element;
// function createElement(tagName: string): Element {
//     return document.createElement(tagName);
// }

// 类型断言
let someValue: unknown = "this is a string";
let strLength: number = (someValue as string).length;

// 非空断言操作符
function processEntity(entity?: { name: string }) {
    // 我们知道 entity 不为 null
    let name = entity!.name;
    console.log("实体名称:", name);
}

// 可选链
interface User2 {
    name?: string;
    address?: {
        street?: string;
        city?: string;
    };
}

function printUserCity(user: User2) {
    console.log("用户城市:", user.address?.city ?? "未知");
}

// 空值合并
function getDisplayName(name: string | null | undefined) {
    return name ?? "匿名用户";
}

console.log("空值合并:", getDisplayName(null));

// 主函数调用
async function main() {
    console.log("TypeScript 高级特性学习");
    console.log("========================");
    
    await asyncExample();
    
    console.log("\n学习完成！");
}

// 运行主函数
main().catch(console.error);