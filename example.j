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
