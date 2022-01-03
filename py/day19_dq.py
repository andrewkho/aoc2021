import time
import heapq
import math
from dataclasses import dataclass, field
from collections import deque, defaultdict, Counter
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

    scanner_pts = []
    with open(f'../inputs/day19/{infile}', 'r') as f:
        pts = None
        for line in f.readlines():
            if len(line.strip()) == 0:
                continue

            m = re.match(r'--- scanner ([0-9]*) ---', line.strip())
            if m:
                pts = list()
                scanner_pts.append(pts)
                continue
            else:
                pts.append([int(x) for x in line.strip().split(',')])

    scanners = []
    for i, pts in enumerate(scanner_pts):
        scanners.append(Scanner(i, np.array(pts, dtype=int)))

    t0 = time.time()
    s0 = scanners[0]
    dq = deque(scanners[1:])
    while len(dq) > 0:
        s1 = dq.popleft()
        found = False
        for rot in ROTS:
            s1.tx = DualQuaternion(
                rot,
                Quaternion.from_translation(np.zeros(3))
            )
            if check_overlap(s0, s1):
                s0.union(s1.transformed_pts())
                found = True
                break
        if not found:
            dq.append(s1)
        else:
            print(len(scanners) - len(dq), '/', len(scanners), len(s0.pts))

    print('1:', len(s0.pts), time.time() - t0, 's')

    t1 = time.time()
    max_l1 = 0
    for i, s0 in enumerate(scanners):
        for s1 in scanners[i+1:]:
            d = np.abs(s0.get_pos() - s1.get_pos()).sum()
            if d > max_l1:
                max_l1 = d

    print('2:', int(np.round(max_l1)), time.time() - t1, 's')


def check_overlap(s0, s1):
    pts0 = s0.pts
    pts1 = s1.transformed_pts()

    dpts = list()
    for j in range(pts1.shape[0]):
        pt1 = pts1[j, :]
        dpts.append(pts0-pt1)

    dptsar = np.concatenate(dpts, axis=0)
    unq, cnts = np.unique(dptsar, axis=0, return_counts=True)
    max_cnt, max_i = max((v, i) for i, v in enumerate(cnts))
    if max_cnt >= 12:
        s1.set_pos(unq[max_i, :])
        return True

    return False


@dataclass
class Quaternion:
    a: np.ndarray

    @classmethod
    def from_rot_mat(cls, m):
        m00, m01, m02 = tuple(m[0, :])
        m10, m11, m12 = tuple(m[1, :])
        m20, m21, m22 = tuple(m[2, :])

        if (m22 < 0):
            if (m00 > m11):
                t = 1 + m00 -m11 -m22
                a = np.array([[t, m01+m10, m20+m02, m12-m21]])
            else:
                t = 1 -m00 + m11 -m22
                a = np.array([[m01+m10, t, m12+m21, m20-m02]])
        else:
            if (m00 < -m11):
                t = 1 -m00 -m11 + m22
                a = np.array([[m20+m02, m12+m21, t, m01-m10]])
            else:
                t = 1 + m00 + m11 + m22
                a = np.array([[m12-m21, m20-m02, m01-m10, t]])
        a *= 0.5 / np.sqrt(t)
        return Quaternion(a)

    @classmethod
    def from_translation(cls, translation: np.ndarray):
        a = np.zeros((1, 4))
        a[:, 1:] = translation
        return Quaternion(a)

    def w(self):
        return self.a[:, 0:1]

    def v(self):
        return self.a[:, 1:]

    def __add__(self, o: 'Quaternion'):
        return Quaternion(self.a + o.a)

    def __sub__(self, o: 'Quaternion'):
        return Quaternion(self.a - o.a)

    def __iadd__(self, o: 'Quaternion'):
        self.a += o.a
        return self

    def __rmul__(self, s: float):
        return Quaternion(self.a * s)

    def __imul__(self, s: float):
        self.a *= s
        return self

    def dot(self, o: 'Quaternion'):
        w = np.multiply(self.w(), o.w())- np.multiply(self.v(), o.v()).sum(axis=1, keepdims=True)
        v = self.w()*o.v() + o.w()*self.v() + np.cross(self.v(), o.v())
        a = np.concatenate([w, v], axis=1)
        return Quaternion(a)

    def conj(self):
        return Quaternion(np.concatenate([self.w(), -1 * self.v()], axis=1))

    def norm(self):
        return (self.dot(self.conj())).w()


@dataclass
class DualQuaternion:
    qr: Quaternion
    qd: Quaternion

    @classmethod
    def from_point(cls, pt: np.ndarray):
        ar = np.zeros((1, 4))
        ar[:, 0] = 1
        ad = np.zeros((1, 4))
        ad[:, 1:] = pt

        return DualQuaternion(Quaternion(ar), Quaternion(ad))

    @classmethod
    def from_points(cls, pts: np.ndarray):
        ar = np.zeros((pts.shape[0], 4))
        ar[:, 0:1] = 1
        ad = np.zeros((pts.shape[0], 4))
        ad[:, 1:] = pts

        return DualQuaternion(Quaternion(ar), Quaternion(ad))

    def set_translation(self, dx: np.ndarray):
        t = Quaternion.from_translation(dx)
        self.qd = 0.5*t.dot(self.qr)

    def get_translation(self):
        return 2*self.qd.dot(self.qr.conj()).a[:, 1:]

    def set_rotation(self, rot: Quaternion):
        translation = self.get_translation()
        self.qr = rot
        self.set_translation(translation)

    def __add__(self, o: 'DualQuaternion'):
        return DualQuaternion(self.qr + o.qr, self.qd + o.qd)

    def __sub__(self, o: 'DualQuaternion'):
        return DualQuaternion(self.qr - o.qr, self.qd - o.qd)

    def __iadd__(self, o: 'DualQuaternion'):
        self.qr += o.qr
        self.qd += o.qd
        return self

    def __rmul__(self, s: float):
        return DualQuaternion(s * self.qr, s * self.qd)

    def dot(self, o: 'DualQuaternion'):
        return DualQuaternion(
            self.qr.dot(o.qr),
            self.qr.dot(o.qd) + self.qd.dot(o.qr)
        )

    def conj(self):
        return DualQuaternion(self.qr.conj(), -1 * self.qd.conj())

    def norm(self):
        return (self.dot(self.conj())).qr.w()


@dataclass
class Scanner:
    n: int
    pts: np.ndarray
    tx: DualQuaternion = field(default_factory=lambda: DualQuaternion.from_point(np.zeros(3)))

    def transformed_pts(self) -> np.ndarray:
        tx = self.tx
        txc = self.tx.conj()
        ptsdq = DualQuaternion.from_points(self.pts)
        return np.round(tx.dot(ptsdq).dot(txc).qd.a[:, 1:]).astype(int)

    def set_pos(self, pos: np.ndarray):
        self.tx.set_translation(pos)

    def get_pos(self):
        return self.tx.get_translation()

    def set_rot(self, rot: Quaternion):
        self.tx.set_rotation(rot)

    def union(self, pts: np.ndarray):
        a = set(tuple(pt) for pt in self.pts)
        a.update(tuple(pt) for pt in pts)
        self.pts = np.array(list(a))


ROTS = []
for x, i in itertools.product(range(3), (-1, 1)):
    xv = np.zeros(3, dtype=float)
    xv[x] = i
    for y, k in itertools.product(range(3), (-1, 1)):
        if y == x:
            continue
        yv = np.zeros(3, dtype=float)
        yv[y] = k
        rotmat = np.stack([xv, yv, np.cross(xv, yv)], axis=0)
        ROTS.append(Quaternion.from_rot_mat(rotmat))


if __name__ == '__main__':
    fire.Fire(main)
