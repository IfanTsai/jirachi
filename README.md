## jirachi

<img align="right" alt="jirachi" src="https://img.caiyifan.cn/typora_picgo/Jirachi.png" />

interpreter for the basic language written in golang

The plan supports the following features:

- [x] Arithmetic Operations (+, -, *, /, ^)
- [x] Comparison Operation (==, !=, >, >=, <, <=)
- [x] Logical Operation (not, and, or)
- [x] Variable
- [x] Judgment Branch Statement (if ... then ... elif ... else ... end)
- [x] Loop Statement (for, while)
- [x] Function
- [x] String
- [x] List
- [ ] Map
- [x] Built-in Functions
- [x] Branch Control Statement (break, continue, return)
- [x] Comment
- [ ] File IO
- [ ] Network IO
- [ ] Coroutine

### start

```bash
make release
sudo make install
# run repl
jirachi
# run jirachi script
./example.j
# or
jirachi example.j
````

### repl

<img src="https://img.caiyifan.cn/typora_picgo/image-20211222234355832.png" alt="image-20211222234355832" style="zoom:80%;" />


### example

```shell
#!/usr/local/bin/jirachi

# loop
fun fib1(n)
    # list comprehension: res = [1, 1, 1]
    res = for i = 0 to 3 then 1

    for i = 3 to n + 1 then
        res = res + (res[i - 1] + res[i - 2])
    end

    return res[n]  # can remove the return keyword
end

# recursion
fib2 = fun(n)      # anonymous function
    if n <= 2 then
        1
    else
        fib2(n - 1) + fib2(n - 2)
    end
end

fib3 = fun(n)
    if n <= 2 then
        return 1
    end

    f1 = 1
    f2 = 1
    f3 = 0

    return for i = 3 to n + 1 then # return final expression value
        f3 = f1 + f2
        f1 = f2
        f2 = f3
    end
end

num2str = fun(number) -> '' + number  # lambda

print_fib_values = fun(fib_fun, n)    # higher order function
    print('fibonacci value is ')

    for i = 1 to n + 1 then
        res = fib_fun(i)
        print(num2str(res) + ' ')
    end

    println('')
end

while true then
    println('please input number:')

    n = input_number()
    print_fib_values(fib1, n)
    print_fib_values(fib2, n)
    print_fib_values(fib3, n)
end
```
