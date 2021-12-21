fib = FUN(n)
    IF n <= 2 THEN
        1
    ELSE
        fib(n - 1) + fib(n - 2)
    END
END

WHILE TRUE THEN
    println('please input number:')
    n = input_number()
    print('fibonacci value is ')
    FOR i = 1 TO n + 1 THEN
        print('' + fib(i) + ' ')
    END
    println('')
END
