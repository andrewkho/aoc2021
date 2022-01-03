import math
from dataclasses import dataclass
from collections import deque
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

    input = []
    with open(f'../inputs/day10/{infile}', 'r') as f:
        for line in f.readlines():
            input.append(list(line.strip()))

    matching = {
        ')': '(',
        ']': '[',
        '>': '<',
        '}': '{',
    }

    points = {
        ')': 3,
        ']': 57,
        '}': 1197,
        '>': 25137,
    }
    points_ac = {
        '(': 1,
        '[': 2,
        '{': 3,
        '<': 4,
    }

    t = 0
    scores = []
    for line in input:
        stack = []
        good = True
        for c in line:
            if c in matching.values():
                stack.append(c)
            else:
                if stack[-1] != matching[c]:
                    t += points[c]
                    good = False
                    break
                else:
                    stack.pop()

        if good:
            score = 0
            for c in stack[::-1]:
                score *= 5
                score += points_ac[c]
            scores.append(score)

    print(f'1: {t}')
    print(f'2: {sorted(scores)[len(scores)//2]}')


if __name__ == '__main__':
    fire.Fire(main)
