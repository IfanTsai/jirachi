#!/usr/local/bin/jirachi

# loop
FUN fib1(n)
    # list comprehension: res = [1, 1, 1]
    res = FOR i = 0 TO 3 THEN 1

    FOR i = 3 TO n + 1 THEN
        res = res + (res[i - 1] + res[i - 2])
    END

    RETURN res[n]  # can remove the RETURN keyword
END

# recursion
fib2 = FUN(n)      # anonymous function
    IF n <= 2 THEN
        1
    ELSE
        fib2(n - 1) + fib2(n - 2)
    END
END

fib3 = FUN(n)
    IF n <= 2 THEN
        RETURN 1
    END

    f1 = 1
    f2 = 1
    f3 = 0

    RETURN FOR i = 3 TO n + 1 THEN  # return final expression value
        f3 = f1 + f2
        f1 = f2
        f2 = f3
    END
END

num2str = FUN(number) -> '' + number  # lambda

print_fib_values = FUN(fib_fun, n)    # higher order function
    print('fibonacci value is ')

    FOR i = 1 TO n + 1 THEN
        res = fib_fun(i)
        print(num2str(res) + ' ')
    END

    println('')
END

WHILE TRUE THEN
    println('please input number:')

    n = input_number()
    print_fib_values(fib1, n)
    print_fib_values(fib2, n)
    print_fib_values(fib3, n)
END
