from dataclasses import dataclass
from typing import *
import numpy as np
from scipy.optimize import linear_sum_assignment
import itertools
import re

import fire


def main(
        #infile: str='test_input.txt',
        infile: str = 'input.txt',
):
    print('hi!')

    input = []
    with open(infile, 'r') as f:
        for line in f.readlines():
            l, r = line.split('|')
            ll = [''.join(sorted(x.strip())) for x in l.split(' ') if len(x) > 0]
            rr = [''.join(sorted(x.strip())) for x in r.split(' ') if len(x) > 0]
            input.append((ll, rr))

    print(input)

    c = 0
    for l, r in input:
        for x in r:
            if len(x) in (2, 3, 4, 7):
                c += 1
    print(f'1: {c}')

    outputs = [
        decode(l, r)
        for l, r in input
    ]

    print(outputs)
    print(f'2: {sum(outputs)}')


possibilities = {
    2: [3, 6],
    3: [1, 3, 6],
    4: [2, 3, 4, 6],
    5: [1, 2, 3, 4, 5, 6, 7],
    6: [1, 2, 3, 4, 5, 6, 7],
    7: [1, 2, 3, 4, 5, 6, 7],
}

lcd = {
    frozenset('25'): 1,
    frozenset('02346'): 2,
    frozenset('02356'): 3,
    frozenset('1235'): 4,
    frozenset('01356'): 5,
    frozenset('013456'): 6,
    frozenset('025'): 7,
    frozenset('0123456'): 8,
    frozenset('012356'): 9,
    frozenset('012456'): 0,
}

letters = {
    'a': 0,
    'b': 1,
    'c': 2,
    'd': 3,
    'e': 4,
    'f': 5,
    'g': 6,
}


def apply_mapping(mapping, ltrs):
    return [
        frozenset(str(mapping[letters[c]])
                  for c in x)
        for x in ltrs
    ]


def decode(l, r):
    A = np.zeros((7, 7), dtype=np.float32)

    for x in l + r:
        ps = possibilities[len(x)]
        for c in x:
            j = letters[c]
            for i in range(7):
                if i+1 not in ps:
                    A[j, i] = np.inf

    def is_feasible(mapping):
        for i, j in enumerate(mapping):
            if A[i, j] > 1.e-3:
                return False
        return True

    # Recurse to try all possibilities
    def recurse(mapping):
        if len(mapping) == 7 and is_feasible(mapping):
            # Check if this is a valid mapping
            if all(x in lcd for x in apply_mapping(mapping, l+r)):
                return mapping
            else:
                return None

        for i in range(7):
            if i in mapping:
                continue
            result = recurse(mapping + [i])
            if result:
                return result

    mapping = recurse([])
    return int(''.join(str(lcd[x]) for x in apply_mapping(mapping, r)))


if __name__ == '__main__':
    fire.Fire(main)
