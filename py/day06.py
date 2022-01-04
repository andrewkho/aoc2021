from dataclasses import dataclass
from typing import *
import numpy as np
import itertools
import re

import fire


def main(
        #infile: str='test_input.txt',
        infile: str = 'input.txt',
):
    print('hi!')

    with open(infile, 'r') as f:
        input = [int(x) for x in f.readlines()[0].split(',')]

    print(len(input))

    h = {
        i: 0 for i in range(9)
    }

    for i in input:
        h[i] += 1

    r = simulate(h.copy(), 80)
    print(r)
    print(f'2: {sum(r.values())}')

    r = simulate(h.copy(), 256)
    print(r)
    print(f'2: {sum(r.values())}')


def simulate(h: Dict[int, int], days: int) -> Dict[int, int]:
    for d in range(days):
        h0 = h[0]
        for i in range(1, 9):
            h[i-1] = h[i]
        h[8] = h0
        h[6] += h0

    return h


if __name__ == '__main__':
    fire.Fire(main)
