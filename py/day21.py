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


def main(infile: str):
    print('hi!')

    lines = []
    with open(infile, 'r') as f:
        for i, line in enumerate(f.readlines()):
            if i == 0:
                p1_0 = int(line.strip()[-2:])
            else:
                p2_0 = int(line.strip()[-2:])

            lines.append(line)

    t0 = time.time()

    state = State((p1_0, p2_0), (0, 0))
    player = 0
    rolls = 0
    while state.score[0] < 1000 and state.score[1] < 1000:
        die = rolls % 100 + 1
        steps = 3*die + 3
        state = state.advance(player, steps)

        rolls += 3
        player = 1 - player

    print('1:', rolls*min(state.score), time.time() - t0, 's')

    t1 = time.time()
    outcomes = defaultdict(lambda: 0)
    for i in itertools.product(range(1, 4), range(1, 4), range(1, 4)):
        outcomes[sum(i)] += 1

    winners = [0, 0]
    player = 0
    hist = defaultdict(lambda: 0, {State(pos=(p1_0, p2_0), score=(0, 0)): 1})
    while len(hist) > 0:
        new_hist = defaultdict(lambda: 0)

        for k, v in hist.items():
            for i, n in outcomes.items():
                state = k.advance(player, i)
                if state.score[player] >= 21:
                    winners[player] += v*n
                else:
                    new_hist[state] += v*n
        hist = new_hist
        player = 1 - player

    print('2:', winners, time.time() - t1, 's', max(winners))


@dataclass(frozen=True)
class State:
    pos: Tuple[int, int]
    score: Tuple[int, int]

    def advance(self, player, steps) -> 'State':
        pos = (self.pos[player]+steps-1) % 10 + 1
        score = self.score[player] + pos
        if player == 0:
            return State(
                (pos, self.pos[1]),
                (score, self.score[1])
            )
        else:
            return State(
                (self.pos[0], pos),
                (self.score[0], score),
            )

if __name__ == '__main__':
    fire.Fire(main)
