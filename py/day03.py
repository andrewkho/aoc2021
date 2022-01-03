import fire


def main(
        #infile: str='test_input.txt',
        infile: str = 'input.txt',
):
    print('hi!')

    with open(infile, 'r') as f:
        input = [
            x.strip()
            for x in f.readlines()
        ]
        width = len(input[0])

    #print(input)

    bits = [0]*width
    for l in input:
        for i, v in enumerate(l):
            bits[i] += int(v)

    g = ['0']*width
    e = ['0']*width
    for i, b in enumerate(bits):
        if b > len(input) / 2:
            g[i] = '1'
        else:
            e[i] = '1'

    print(f"1: {int(''.join(g), 2) * int(''.join(e), 2)}")

    # first loop
    input2 = input.copy()
    i = 0
    while True:
        b = sum([int(x[i]) for x in input2])
        if b >= len(input2) / 2:
            c = '1'
        else:
            c = '0'
        input2 = [x for x in input2 if x[i] == c]
        if len(input2) <= 1:
            break
        i += 1

    o2 = input2[0]

    input2 = input.copy()
    i = 0
    while True:
        b = sum([int(x[i]) for x in input2])
        if b >= len(input2) / 2:
            c = '0'
        else:
            c = '1'
        input2 = [x for x in input2 if x[i] == c]
        if len(input2) <= 1:
            break
        i += 1

    co2 = input2[0]
    print(o2, co2)
    print(f"2: {int(''.join(o2), 2) * int(''.join(co2), 2)}")


if __name__ == '__main__':
    fire.Fire(main)
