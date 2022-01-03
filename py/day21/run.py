import time
import heapq
import math
from dataclasses import dataclass, field
from collections import deque, defaultdict, Counter
from typing import *
from scipy.signal import convolve2d

import numpy as np
import itertools
import re

import fire


def main(
        #infile: str='test_input.txt',
        infile: str = 'input.txt',
):
    print('hi!')

    lines = []
    with open(f'../inputs/day21/{infile}', 'r') as f:
        for i, line in enumerate(f.readlines()):
            if i == 0:
                p1_0 = int(line.strip()[-2:])
            else:
                p2_0 = int(line.strip()[-2:])

            lines.append(line)

    def advance(pos, steps):
        pos += steps
        if pos > 10:
            pos = (pos - 1) % 10 + 1
        return pos

    t0 = time.time()

    s1, s2 = 0, 0
    p1 = p1_0
    p2 = p2_0
    turns = 0
    die = 0
    while s1 < 1000 and s2 < 1000:
        steps = 0
        for i in range(3):
            die = die+1 if die < 100 else 1
            steps += die

        if turns % 2 == 0:
            p1 = advance(p1, steps)
            s1 += p1
        else:
            p2 = advance(p2, steps)
            s2 += p2

        turns += 1

    print('1:', turns, s1, 3*turns*s1, s2, 3*turns*s2, time.time() - t0, 's')

    t1 = time.time()
    outcomes = defaultdict(lambda: 0)
    for i in itertools.product(range(1, 4), range(1, 4), range(1, 4)):
        outcomes[sum(i)] += 1

    h1_winners = 0
    h2_winners = 0
    hist = defaultdict(lambda: 0, {State(p1_0, 0, p2_0, 0): 1})
    while len(hist) > 0:
        new_hist = defaultdict(lambda: 0)

        for k, v in hist.items():
            for i, n in outcomes.items():
                pos1 = advance(k.pos1, i)
                scr1 = k.score1 + pos1
                if scr1 >= 21:
                    h1_winners += v*n
                else:
                    for j, m in outcomes.items():
                        pos2 = advance(k.pos2, j)
                        scr2 = k.score2 + pos2
                        if scr2 >= 21:
                            h2_winners += v*n*m
                        else:
                            new_hist[State(pos1, scr1, pos2, scr2)] += v*n*m

        hist = new_hist

    print('2:', h1_winners, h2_winners, time.time() - t1, 's', max([h1_winners, h2_winners]))


@dataclass(frozen=True)
class State:
    pos1: int
    score1: int
    pos2: int
    score2: int


if __name__ == '__main__':
    fire.Fire(main)
