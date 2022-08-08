函数模板

对比函数重载(同一作用域内函数名相同，参数列表不同的函数)
函数模板只需要一个函数就实现了函数重载的部分功能
(参数个数相同类型不同，函数重载需要定义多个同名参数列表不同的函数)

template<typename T, typename Y>  // 这也可以写 template<class T, class Y> 此处的 class 和 typename 作用相同
void tfunc(T& t, Y& y) {
cout << t << " " << y << endl;
}

int n = 2;
double d = 2.1;
tfunc(n, d);

// 运行结果：2 2.1



函数模板具体化，函数模板具体化就是将某一(某几)个要处理的类型单独处理
需要单独写一个实现，形式是 template<> void fun(type& t);
函数模板的具体化和普通函数可以同时存在，调用顺序是：
普通函数 > 函数模板具体化 > 模板函数

// ====== 测试一：函数模板针对特殊数据类型具体化 ======

// 结构体
struct Node {
int val;
Node* next;
};

// 函数模板
template<typename T>
void tfunc(const T& t) {
cout << "template: " << t << endl;
}

// 函数模板具体化(用于处理Node类型数据)
template<>
void tfunc<Node>(const Node& node) {
cout << "template<Node>: " << node.val << endl;
}

// 函数模板具体化(用于处理int类型数据)
template<>
void tfunc<int>(const int& n) {
cout << "template<int>: " << n << endl;
}

// 普通函数
void tfunc(const int& n) {
cout << "tfunc(): " << n << endl;
}

double d = 2.1;
tfunc(d);    // 函数模板未具体化double类型函数，调用模板
Node node{ 2, nullptr };
tfunc(node); // 函数模板具体化Node类型函数，调用函数模板的具体化
int n = 2;
tfunc(n);    // 函数模板具体化int类型函数，也存在普通函数，调用普通函数


// ====== 测试二：函数模板部分具体化 ======
template<typename T1, typename T2>
void tfunc(T1 t1, T2 t2) {
cout << typeid(T1).name() << " and " << typeid(T2).name() <<": " << t1 << " " << t2 << endl;
}

template<typename T1>
void tfunc(T1 t1, int i) {
cout << typeid(T1).name() << " and " << "int: " << t1 << " " << i << endl;
}

template<typename T2>
void tfunc(long l, T2 t2) {
cout << "long and " << typeid(T2).name() << ": " << l << " " << t2 << endl;
}

template<>
void tfunc(short l, int i) {
cout << "long and int: " << l << " " << i << endl;
}

// 分别调用以上四个模板函数
tfunc(char('c'), char('c'));
tfunc(char('c'), int(10));
tfunc(long(10), char('c'));
tfunc(short(10), int(10));


函数模板实例化，让编译器生成指定类型的函数定义，不用写函数的实现，形式是 template void fun(type& t);



// 函数模板
template<typename T>
void tfunc(const T& t) {
cout << "template: " << t << endl;
}


// 函数模板实例化，不用写函数的实现，编译器会生成该类型的模板具体化函数
template void tfunc<char>(const char& c);

类模板

类模板可以指定默认模板参数(函数模板不可以)，跟函数参数的默认值一样，必须从右向左连续赋值默认类型，
如果实例化对象时又传递了类型，默认类型会被覆盖掉，跟函数参数是一样的
创建对象时需要传递模板参数列表，模板参数列表加在类名后面 ClassName< typename T > classN;
如果类的模板参数列表有默认值，可以不传模板参数，但一定要加 <> 如 ClassName< > classN;
创建堆区对象的时候，所有的类名称后面都要加模板参数列表，如 ClassName< typename T >* classN = new ClassName< typename T>; 
除了类内，其他地方出现 ClassName 的地方一般都要加模板参数列表

template<typename T = int, typename Y = char> // 此处指定了模板默认参数，部分指定必须从右到左指定
class Test {
public:
    Test(T t, Y y) : t(t), y(y) { // 这里使用了初始化列表

    }

    void tfunc();

private:
    T t;
    Y y;
};

template<typename T, typename Y> // 类模板的函数在类外实现，需要加上模板参数列表，但数不需要加指定的默认模板参
void Test<T, Y>::tfunc() { // 类外使用Test需要加模板参数
cout << t << " " << y << endl;
}

int n = 2;
double d = 2.1;
Test<int, double> test(n, d); // 此处如果使用默认模板参数可定义为 Test<> test(int(2), char('a'));
test.tfunc();

// 运行结果：2 2.1

类模板的继承，类模板被继承后参数的传递方式主要有两种，一种是直接在子类继承父类的时候，为父类指定固定的类型；
二是通过子类模板参数列表传递

// ====== 测试一 ======
template<typename T, typename Y>
class A {
public:
A(T t, Y y) {

    }
};

class Test : public A<int, double> { // 父类是类模板，子类是普通类
public:
Test() : A<int, double>(2, 2.1) {

    }
};

Test();

// ====== 测试二 ====== 二是通过子类模板参数列表传递
template<typename T, typename Y>
class A {
public:
A(T t) {

    }
};

template<typename X, typename Z, typename P>
class Test : public A<X, P> {
public:
Test(X x, Z z, P p) : A<X, P>(x) {

    }
};

Test<int, double, char>(int(2), double(2.1), char('a'));

类模板的多态，在创建对象时，分为：
子类没有模板(Father<char, short>* cf = new Son)
子类有模板(Father<short, char>* cf = new Son<short, int, char>) 两种，子类和父类的模板参数列表可以不一样，但一定要对应好
// ====== 测试一 ======
template<typename T, typename Y>
class A {
public:
virtual void tfunc(T t, Y y) = 0; // 声明一个虚函数
};

class Test : public A<int, double> {
public:
virtual void tfunc(int n, double d) {
cout << n << " " << d << endl;
}
};

// 父类是类模板，子类是普通类，在多态情况下父类需要指定模板参数，子类就不用了
A<int, double>* a = new Test;
a->tfunc(2, 2.1);
// 运行结果：2 2.1

// ====== 测试二 ======
template<typename T, typename Y>
class A {
public:
virtual void tfunc(T t, Y y) = 0;
};

template<typename X, typename Z, typename P>
class Test : public A<X, P> {
public:
virtual void tfunc(X x, P p) {
cout << x << " " << p << endl;
}
};

// 父类是类模板，子类是类模板，在多态情况下父类和子类都需要指定模板参数
A<int, double>* a = new Test<int, char, double>;
a->tfunc(2, 2.1);
// 运行结果：2 2.1

类模板具体化，类模板的具体化分为部分具体化和全部具体化两种
template<typename T1, typename T2>
class Test {
public:
Test() {
cout << "T1 and T2" << endl;
}
};

// 部分具体化
template<typename T1>
class Test<T1, int> {
public:
Test() {
cout << "T1 and int" << endl;
}
};

// 部分具体化
template<typename T2>
class Test<long, T2> {
public:
Test() {
cout << "long and T2" << endl;
}
};

// 全部具体化
template<>
class Test<long, int> {
public:
Test() {
cout << "long and int" << endl;
}
};

// 分别创建上面四个类
Test<char, char>();
Test<char, int>();
Test<long, char>();
Test<long, int>();

成员模板

成员模板简单说就是模板中的模板
class Base1 { };
class Base2 { };

class Test1 : public Base1 { };
class Test2 : public Base2 { };

template<typename T1, typename T2>
class Pair {
public:
T1 t1;
T2 t2;

    Pair(T1 t1, T2 t2) : t1(t1), t2(t2) {

    }

    // 类模板中的成员模板
    template<typename U1, typename U2>
    Pair(const Pair<U1, U2>& pair) : t1(pair.t1), t2(pair.t2){

    }
};

Pair<Base1*, Base2*>(Pair<Test1*, Test2*>(new Test1, new Test2));