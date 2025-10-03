// TypeScript 引用类型学习文件
// 演示 TypeScript 中的引用类型、对象引用、内存管理等概念

// 基础类型定义
interface RefPerson {
    name: string;
    age: number;
    city: string;
}

interface RefAddress {
    street: string;
    city: string;
    zipCode: string;
}

// 类定义
class RefPersonClass {
    constructor(
        public name: string,
        public age: number,
        public address: RefAddress
    ) {}

    getInfo(): string {
        return `${this.name}, ${this.age}岁, 住址: ${this.address.city}`;
    }

    updateAge(newAge: number): void {
        this.age = newAge;
    }

    updateAddress(newAddress: RefAddress): void {
        this.address = newAddress;
    }
}

// 导出一个空对象使此文件成为模块，避免全局作用域冲突
export {};

// 基本引用类型演示
function basicReferenceDemo(): void {
    console.log("\n=== 基本引用类型演示 ===");
    
    // 原始类型 vs 引用类型
    console.log("--- 原始类型 vs 引用类型 ---");
    
    // 原始类型（值类型）
    let a = 10;
    let b = a;  // 复制值
    b = 20;
    console.log(`原始类型 - a: ${a}, b: ${b}`); // a: 10, b: 20
    
    // 引用类型
    let obj1 = { value: 10 };
    let obj2 = obj1;  // 复制引用
    obj2.value = 20;
    console.log(`引用类型 - obj1.value: ${obj1.value}, obj2.value: ${obj2.value}`); // 都是 20
    
    // 对象引用比较
    console.log("--- 引用比较 ---");
    let person1: RefPerson = { name: "张三", age: 25, city: "北京" };
    let person2: RefPerson = { name: "张三", age: 25, city: "北京" };
    let person3 = person1;
    
    console.log(`person1 === person2: ${person1 === person2}`); // false (不同引用)
    console.log(`person1 === person3: ${person1 === person3}`); // true (相同引用)
    
    // 深度比较
    function deepEqual(obj1: any, obj2: any): boolean {
        return JSON.stringify(obj1) === JSON.stringify(obj2);
    }
    
    console.log(`深度比较 person1 和 person2: ${deepEqual(person1, person2)}`); // true
}

// 数组引用演示
function arrayReferenceDemo(): void {
    console.log("\n=== 数组引用演示 ===");
    
    // 数组引用
    let arr1 = [1, 2, 3, 4, 5];
    let arr2 = arr1;  // 引用复制
    
    console.log("原始数组:", arr1);
    
    arr2.push(6);
    console.log("修改 arr2 后:");
    console.log("arr1:", arr1);
    console.log("arr2:", arr2);
    
    // 数组浅拷贝
    let arr3 = [...arr1];  // 展开运算符
    let arr4 = Array.from(arr1);  // Array.from
    let arr5 = arr1.slice();  // slice 方法
    
    arr3.push(7);
    console.log("浅拷贝后修改 arr3:");
    console.log("arr1:", arr1);
    console.log("arr3:", arr3);
    
    // 嵌套数组的引用问题
    let nestedArr1 = [[1, 2], [3, 4]];
    let nestedArr2 = [...nestedArr1];  // 浅拷贝
    
    nestedArr2[0].push(3);  // 修改嵌套数组
    console.log("嵌套数组浅拷贝问题:");
    console.log("nestedArr1:", nestedArr1);
    console.log("nestedArr2:", nestedArr2);
    
    // 深拷贝解决方案
    let nestedArr3 = JSON.parse(JSON.stringify(nestedArr1));
    nestedArr3[0].push(4);
    console.log("深拷贝后:");
    console.log("nestedArr1:", nestedArr1);
    console.log("nestedArr3:", nestedArr3);
}

// 对象引用和修改
function objectReferenceDemo(): void {
    console.log("\n=== 对象引用和修改演示 ===");
    
    // 对象属性修改
    let person: RefPerson = { name: "李四", age: 30, city: "上海" };
    let personRef = person;
    
    console.log("原始对象:", person);
    
    // 通过引用修改对象
    personRef.age = 31;
    personRef.city = "深圳";
    
    console.log("通过引用修改后:");
    console.log("person:", person);
    console.log("personRef:", personRef);
    
    // 对象浅拷贝-展开运算符和Object.assign
    let personCopy1 = { ...person };  // 展开运算符
    let personCopy2 = Object.assign({}, person);  // Object.assign
    
    personCopy1.age = 32;
    console.log("浅拷贝修改后:");
    console.log("person:", person);
    console.log("personCopy1:", personCopy1);
    personCopy2.age = 33;
    console.log("assign拷贝修改后:");
    console.log("person:", person);
    console.log("personCopy2:", personCopy2);
    
    // 嵌套对象引用问题
    interface RefPersonWithAddress {
        name: string;
        age: number;
        address: RefAddress;
    }
    
    let personWithAddr: RefPersonWithAddress = {
        name: "王五",
        age: 28,
        address: { street: "中山路123号", city: "广州", zipCode: "510000" }
    };
    
    let personWithAddrCopy = { ...personWithAddr };  // 浅拷贝
    personWithAddrCopy.address.city = "杭州";  // 修改嵌套对象
    
    console.log("嵌套对象浅拷贝问题:");
    console.log("原对象地址:", personWithAddr.address.city);
    console.log("拷贝对象地址:", personWithAddrCopy.address.city);
}

// 函数参数引用
function functionParameterDemo(): void {
    console.log("\n=== 函数参数引用演示 ===");
    
    // 原始类型参数（值传递）
    function modifyPrimitive(x: number): void {
        x = x + 10;
        console.log("函数内部 x:", x);
    }
    
    let num = 5;
    console.log("调用前 num:", num);
    modifyPrimitive(num);
    console.log("调用后 num:", num);
    
    // 对象参数（引用传递）
    function modifyObject(person: RefPerson): void {
        person.age = person.age + 1;
        person.city = "修改后的城市";
        console.log("函数内部 person:", person);
    }
    
    let person: RefPerson = { name: "赵六", age: 25, city: "成都" };
    console.log("调用前 person:", person);
    modifyObject(person);
    console.log("调用后 person:", person);
    
    // 数组参数（引用传递）
    function modifyArray(arr: number[]): void {
        arr.push(100);
        arr[0] = 999;
        console.log("函数内部 arr:", arr);
    }
    
    let numbers = [1, 2, 3, 4, 5];
    console.log("调用前 numbers:", numbers);
    modifyArray(numbers);
    console.log("调用后 numbers:", numbers);
    
    // 避免意外修改的方法
    function safeModifyArray(arr: readonly number[]): number[] {
        // 创建新数组而不是修改原数组
        return [...arr, 100];
    }
    
    let safeNumbers = [1, 2, 3];
    let newNumbers = safeModifyArray(safeNumbers);
    console.log("安全修改:");
    console.log("原数组:", safeNumbers);
    console.log("新数组:", newNumbers);
}

// 类实例引用
function classInstanceDemo(): void {
    console.log("\n=== 类实例引用演示 ===");
    
    let address1: RefAddress = { street: "解放路456号", city: "武汉", zipCode: "430000" };
    let person1 = new RefPersonClass("孙七", 35, address1);
    let person2 = person1;  // 引用复制
    
    console.log("原始实例:", person1.getInfo());
    
    // 通过引用修改实例
    person2.updateAge(36);
    console.log("通过引用修改年龄后:");
    console.log("person1:", person1.getInfo());
    console.log("person2:", person2.getInfo());
    
    // 修改嵌套对象
    let newAddress: RefAddress = { street: "人民路789号", city: "西安", zipCode: "710000" };
    person2.updateAddress(newAddress);
    console.log("修改地址后:");
    console.log("person1:", person1.getInfo());
    
    // 类实例的浅拷贝
    let person3 = Object.assign(Object.create(Object.getPrototypeOf(person1)), person1);
    person3.updateAge(37);
    console.log("浅拷贝修改后:");
    console.log("person1 年龄:", person1.age);
    console.log("person3 年龄:", person3.age);
}

// 闭包和引用
function closureReferenceDemo(): void {
    console.log("\n=== 闭包和引用演示 ===");
    
    // 闭包捕获引用
    function createCounter() {
        let count = 0;
        return {
            increment: () => ++count,
            decrement: () => --count,
            getValue: () => count
        };
    }
    
    let counter1 = createCounter();
    let counter2 = createCounter();
    
    console.log("counter1 初始值:", counter1.getValue());
    console.log("counter2 初始值:", counter2.getValue());
    
    counter1.increment();
    counter1.increment();
    counter2.increment();
    
    console.log("操作后:");
    console.log("counter1 值:", counter1.getValue());
    console.log("counter2 值:", counter2.getValue());
    
    // 闭包捕获对象引用
    function createPersonManager(person: RefPerson) {
        return {
            getAge: () => person.age,
            setAge: (age: number) => { person.age = age; },
            getPerson: () => person
        };
    }
    
    let person: RefPerson = { name: "周八", age: 40, city: "重庆" };
    let manager = createPersonManager(person);
    
    console.log("闭包捕获对象:");
    console.log("初始年龄:", manager.getAge());
    
    manager.setAge(41);
    console.log("通过闭包修改后:");
    console.log("manager 获取的年龄:", manager.getAge());
    console.log("原对象年龄:", person.age);
}

// WeakMap 和 WeakSet 演示
function weakReferencesDemo(): void {
    console.log("\n=== 弱引用演示 ===");
    
    // WeakMap 演示
    let wm = new WeakMap();
    let obj1 = { name: "对象1" };
    let obj2 = { name: "对象2" };
    
    wm.set(obj1, "obj1 的值");
    wm.set(obj2, "obj2 的值");
    
    console.log("WeakMap 中 obj1 的值:", wm.get(obj1));
    console.log("WeakMap 中 obj2 的值:", wm.get(obj2));
    
    // WeakSet 演示
    let ws = new WeakSet();
    ws.add(obj1);
    ws.add(obj2);
    
    console.log("WeakSet 包含 obj1:", ws.has(obj1));
    console.log("WeakSet 包含 obj2:", ws.has(obj2));
    
    // 注意：当对象被垃圾回收时，WeakMap 和 WeakSet 中的引用也会被自动清除
    console.log("WeakMap 和 WeakSet 不会阻止对象被垃圾回收");
}

// 内存泄漏预防
function memoryLeakPreventionDemo(): void {
    console.log("\n=== 内存泄漏预防演示 ===");
    
    // 事件监听器内存泄漏示例
    class EventEmitter {
        private listeners: { [event: string]: Function[] } = {};
        
        on(event: string, callback: Function): void {
            if (!this.listeners[event]) {
                this.listeners[event] = [];
            }
            this.listeners[event].push(callback);
        }
        
        off(event: string, callback: Function): void {
            if (this.listeners[event]) {
                this.listeners[event] = this.listeners[event].filter(cb => cb !== callback);
            }
        }
        
        emit(event: string, data?: any): void {
            if (this.listeners[event]) {
                this.listeners[event].forEach(callback => callback(data));
            }
        }
        
        removeAllListeners(): void {
            this.listeners = {};
        }
    }
    
    let emitter = new EventEmitter();
    let obj = { name: "监听对象" };
    
    // 添加监听器
    let callback = (data: any) => {
        console.log(`${obj.name} 收到事件:`, data);
    };
    
    emitter.on("test", callback);
    emitter.emit("test", "测试数据");
    
    // 清理监听器以防止内存泄漏
    emitter.off("test", callback);
    console.log("已清理事件监听器");
    
    // 定时器内存泄漏预防
    let timerId = setTimeout(() => {
        console.log("定时器执行");
    }, 1000);
    
    // 清理定时器
    clearTimeout(timerId);
    console.log("已清理定时器");
}

// 深拷贝实现
function deepCopyDemo(): void {
    console.log("\n=== 深拷贝实现演示 ===");
    
    // 简单深拷贝（JSON 方法）
    function simpleDeepCopy<T>(obj: T): T {
        return JSON.parse(JSON.stringify(obj));
    }
    
    // 完整深拷贝实现
    function deepCopy<T>(obj: T, visited = new WeakMap()): T {
        // 处理 null 和 undefined
        if (obj === null || typeof obj !== "object") {
            return obj;
        }
        
        // 处理循环引用
        if (visited.has(obj as any)) {
            return visited.get(obj as any);
        }
        
        // 处理日期
        if (obj instanceof Date) {
            return new Date(obj.getTime()) as any;
        }
        
        // 处理数组
        if (Array.isArray(obj)) {
            const arrCopy: any[] = [];
            visited.set(obj as any, arrCopy);
            for (let i = 0; i < obj.length; i++) {
                arrCopy[i] = deepCopy(obj[i], visited);
            }
            return arrCopy as any;
        }
        
        // 处理对象
        const objCopy = {} as any;
        visited.set(obj as any, objCopy);
        
        for (const key in obj) {
            if (obj.hasOwnProperty(key)) {
                objCopy[key] = deepCopy(obj[key], visited);
            }
        }
        
        return objCopy;
    }
    
    // 测试深拷贝
    let complexObj = {
        name: "复杂对象",
        age: 30,
        hobbies: ["读书", "游泳"],
        address: {
            city: "北京",
            details: {
                street: "朝阳路",
                number: 123
            }
        },
        date: new Date()
    };
    
    // 添加循环引用
    (complexObj as any).self = complexObj;
    
    console.log("原始对象:", complexObj.name);
    
    let deepCopied = deepCopy(complexObj);
    deepCopied.name = "深拷贝对象";
    deepCopied.address.city = "上海";
    deepCopied.hobbies.push("跑步");
    
    console.log("深拷贝后:");
    console.log("原对象名称:", complexObj.name);
    console.log("原对象城市:", complexObj.address.city);
    console.log("原对象爱好:", complexObj.hobbies);
    console.log("拷贝对象名称:", deepCopied.name);
    console.log("拷贝对象城市:", deepCopied.address.city);
    console.log("拷贝对象爱好:", deepCopied.hobbies);
}

// 泛型引用类型
function genericReferenceDemo(): void {
    console.log("\n=== 泛型引用类型演示 ===");
    
    // 泛型容器类
    class Container<T> {
        private items: T[] = [];
        
        add(item: T): void {
            this.items.push(item);
        }
        
        get(index: number): T | undefined {
            return this.items[index];
        }
        
        getAll(): T[] {
            return [...this.items]; // 返回副本
        }
        
        getAllRef(): T[] {
            return this.items; // 返回引用
        }
        
        size(): number {
            return this.items.length;
        }
    }
    
    // 使用泛型容器
    let personContainer = new Container<RefPerson>();
    personContainer.add({ name: "泛型人员1", age: 25, city: "北京" });
    personContainer.add({ name: "泛型人员2", age: 30, city: "上海" });
    
    console.log("容器大小:", personContainer.size());
    
    let allPersons = personContainer.getAll(); // 获取副本
    let allPersonsRef = personContainer.getAllRef(); // 获取引用
    
    allPersons[0].name = "修改副本";
    allPersonsRef[1].name = "修改引用";
    
    console.log("修改后:");
    console.log("原容器中的第一个人员:", personContainer.get(0)?.name);
    console.log("原容器中的第二个人员:", personContainer.get(1)?.name);
    
    // 泛型函数引用
    function processItems<T>(items: T[], processor: (item: T) => T): T[] {
        return items.map(processor);
    }
    
    let numbers = [1, 2, 3, 4, 5];
    let doubled = processItems(numbers, x => x * 2);
    
    console.log("原数组:", numbers);
    console.log("处理后:", doubled);
}

// 异步引用处理
async function asyncReferenceDemo(): Promise<void> {
    console.log("\n=== 异步引用处理演示 ===");
    
    // 异步操作中的引用
    let sharedData = { count: 0, status: "初始" };
    
    async function asyncOperation1(data: { count: number; status: string }): Promise<void> {
        return new Promise(resolve => {
            setTimeout(() => {
                data.count += 1;
                data.status = "操作1完成";
                console.log("异步操作1完成:", data);
                resolve();
            }, 100);
        });
    }
    
    async function asyncOperation2(data: { count: number; status: string }): Promise<void> {
        return new Promise(resolve => {
            setTimeout(() => {
                data.count += 2;
                data.status = "操作2完成";
                console.log("异步操作2完成:", data);
                resolve();
            }, 50);
        });
    }
    
    console.log("开始异步操作, 初始数据:", sharedData);
    
    // 并发执行异步操作
    await Promise.all([
        asyncOperation1(sharedData),
        asyncOperation2(sharedData)
    ]);
    
    console.log("所有异步操作完成:", sharedData);
    
    // Promise 链中的引用
    let chainData = { value: 1 };
    
    await Promise.resolve(chainData)
        .then(data => {
            data.value *= 2;
            console.log("Promise 链步骤1:", data);
            return data;
        })
        .then(data => {
            data.value += 10;
            console.log("Promise 链步骤2:", data);
            return data;
        });
    
    console.log("Promise 链完成后的数据:", chainData);
}

// 主函数
async function main(): Promise<void> {
    console.log("=== TypeScript 引用类型学习 ===");
    
    // 1. 基本引用类型
    basicReferenceDemo();
    
    // 2. 数组引用
    arrayReferenceDemo();
    
    // 3. 对象引用
    objectReferenceDemo();
    
    // 4. 函数参数引用
    functionParameterDemo();
    
    // 5. 类实例引用
    classInstanceDemo();
    
    // 6. 闭包和引用
    closureReferenceDemo();
    
    // 7. 弱引用
    weakReferencesDemo();
    
    // 8. 内存泄漏预防
    memoryLeakPreventionDemo();
    
    // 9. 深拷贝
    deepCopyDemo();
    
    // 10. 泛型引用
    genericReferenceDemo();
    
    // 11. 异步引用处理
    await asyncReferenceDemo();
    
    console.log("\n=== TypeScript 引用类型总结 ===");
    console.log("1. JavaScript/TypeScript 中对象、数组、函数都是引用类型");
    console.log("2. 引用类型的赋值是复制引用，不是复制值");
    console.log("3. 修改引用类型会影响所有指向同一对象的变量");
    console.log("4. 使用浅拷贝和深拷贝来避免意外的引用修改");
    console.log("5. 注意闭包中的引用捕获");
    console.log("6. 使用 WeakMap 和 WeakSet 来避免内存泄漏");
    console.log("7. 在异步操作中要特别注意共享引用的修改");
    
    console.log("\n学习完成！");
}

// 执行主函数
main().catch(console.error);