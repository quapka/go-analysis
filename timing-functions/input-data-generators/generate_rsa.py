from Crypto.PublicKey import RSA
import gmpy
import argparse
from random import choice, getrandbits

LOW_WEIGHT = 1
HIGH_WEIGHT = 2
RANDOM_WEIGHT = 3
WEIGHTS = {
    'l' : LOW_WEIGHT,
    'h' : HIGH_WEIGHT,
    'r' : RANDOM_WEIGHT,
}

def weight_name(weight):
    for key, value in WEIGHTS.items():
        if value == weight:
            return key

def generate_rsa(iterations=1000, size=1024, weight=RANDOM_WEIGHT):
    key = RSA.generate(size)
    phi_n = (key.p - 1) * (key.q - 1)
    rand = getrandbits(size - 1)
    ciphertext = getrandbits(size - 1)

    out_filename = 'rsa_size-{}_weight-{}.csv'.format(size, weight_name(weight))
    with open(out_filename, 'w') as f:
        for it in range(iterations):
            new_e = 0
            while new_e == 0:
                new_d = generate_private(weight=weight, modulus=phi_n)
                new_e = int(gmpy.invert(new_d, phi_n))

            # check, that the inverse works
            assert ( (new_d * new_e) % phi_n == 1 )
            # NOTE be aware, overriding .e and .d does not change it internally!
            key.e = new_e
            key.d = new_d

            assert( new_e == int(new_e)) 

            # TODO finish the row and write it to the file
            row = 'col_id-{it}_w-{weight};{rand:X};{mod:X};{e:X};{d:X};{p:X};{q:X};{cipher:X}\n'.format(
                    it=it,
                    weight=sum([1 if x == '1' else 0 for x in bin(key.d)]),
                    rand=rand,
                    mod=key.n,
                    e=new_e,
                    d=new_d,
                    p=key.p,
                    q=key.q,
                    cipher=ciphertext,
                    )
            f.write(row)
    return key


def generate_private(weight=RANDOM_WEIGHT, modulus=2048):
    bit_len = len(bin(modulus)) - 3 # minus 2 for the 0b prefix
    if weight == LOW_WEIGHT:
        d = 1 << bit_len
        for _ in range(choice(range(1, 5))):
            d |= 1 << choice(range(bit_len))
        d |= 1
    elif weight == HIGH_WEIGHT:
        d = int('1' * bit_len, 2)
        for _ in range(choice(range(1, 5))):
            d = d ^ (1 << choice(range(bit_len)))
        d |= 1
    elif weight == RANDOM_WEIGHT:
        d = choice(range(modulus))
        d |= 1
    else:
        raise Exception("Not implemented weight: %s" % weight)
    if d >= modulus:
        return 0

    return d


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Generate a csv file with RSA values')
    parser.add_argument('-i', '--iterations', type=int, required=True)

    parser.add_argument('-w', '--weight', choices=WEIGHTS.keys(), required=True)
    parser.add_argument('-s', '--size', type=int, required=False, default=1024)

    args = parser.parse_args()


    generate_rsa(iterations=args.iterations,
                 weight=WEIGHTS[args.weight],
                 size=args.size)
