from dataclasses import dataclass
from typing import *
import numpy as np
import itertools
import re
import time

import fire


def main(
        #infile: str='test_input.txt',
        infile: str = 'input.txt',
):
    print('hi!')

    with open(infile, 'r') as f:
        input = [int(x) for x in f.readlines()[0].split(',')]

    input = sorted(input)

    def cost(x):
        return sum(abs(i-x) for i in input)

    t0 = time.time()
    med = np.median(input)
    print(f'1: {med}, {cost(med)}', time.time() - t0, 's')

    t1 = time.time()
    med2 = binsearch(input, cost)
    print(f'1: {med2}, {cost(med2)}', time.time() - t1, 's')

    def cost2(x):
        # sum(range(N)) = N*(N+1)/2
        return sum(
            abs(x-i)*(abs(x-i)+1)/2 for i in input
        )

    mincost = cost2(binsearch(input, cost2))
    print(f'mincost2: {mincost}')


def binsearch(input, cost):
    l = input[0]
    r = input[-1]

    while l <= r:
        m = (l+r) // 2
        cl = cost(m-1)
        cm = cost(m)
        if cm < cl:
            l = m+1
        else:
            r = m-1

    return l-1


if __name__ == '__main__':
    fire.Fire(main)
