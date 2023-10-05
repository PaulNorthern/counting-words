# counting-words
Parallel reading of a file in N gorutins and accumulation of the number of words passed in the program arguments.
For example, we are looking for 5 words in the file: pleasure explain itself mistaken pain

In the first case, 2 goroutines were involved to process the file, which took 44ms.

![Снимок экрана 2023-10-05 204659](https://github.com/PaulNorthern/counting-words/assets/24978186/940fce3a-486d-4a48-9e2a-92e17e2138a4)

In the second case, I set 32 goroutines in the countThreads variable and it took 23ms.

![image](https://github.com/PaulNorthern/counting-words/assets/24978186/77a3a6b5-e109-4a3b-b446-7d31d4d3b9b2)

This gave an idea of how keyword counting works in the same MS Word for hundreds of pages of text.
