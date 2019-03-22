#Functions to analyze:

crypto/elliptic/ScalarMult
crypto/elliptic/p256.go: p256ScalarMult
crypto/elliptic/p256.go: p256Scalar4
crypto/elliptic/p256.go: p256Scalar8
crypto/elliptic/p256_asm.go: p256PointAddAsm
crypto/elliptic/p256_asm.go: p256PointDoubleAsm
crypto/elliptic/p256_asm.go: Inverse
crypto/elliptic/p256_asm.go: p256StorePoint
crypto/elliptic/p256_asm.go: maybeReduceModP
crypto/rsa.DecryptPKCS1v15SessionKey
crypto/rsa.decryptPKCS1v15
crypto/rsa.SignPKCS1v15

*Note: for ECC, we only work with NISTP-256.

Possible bonus:
crypto/rsa.leftPad
(*math/big.Int).ModInverse
crypto/elliptic.p256FromMont
(*crypto/elliptic.CurveParams).doubleJacobian

##ScalarMult
inputs: Bx, By *big.Int, k []byte

Ideas: try different scalars (they have to be between 1 and 2^256−432420386565659656852420866394968145599 (the order of NISTP-256)), choose them with low HW/ high HW/ randomly. Apply the same pattern to Bx and then compute By so that (Bx, By) lies on the curve.

##p256ScalarMult
inputs: bigX, bigY *big.Int, scalar []byte
Ideas: Exactly the same as above, try to find a difference (unless this function is called by the above one).


##p256Scalar4
inputs: out *[p256Limbs]uint32

Ideas: try all possible values.

##p256Scalar8
inputs: out *[p256Limbs]uint32

Ideas: try all possible values.

##p256PointAddAsm
inputs: res, in1, in2 []uint64

Ideas: try low HW/high HW/random values for in1, in2.

##p256PointAddAsm
inputs: res, in []uint64

Ideas: try low HW/high HW/random values for in.

##Inverse
inputs: k *big.Int

Ideas: try different k's around the range of the order of NISTP-256.

##p256StorePoint
inputs: r *[16 * 4 * 3]uint64, index int

Ideas: this is just a simple function, try low HW/high HW/random values.

##maybeReduceModP
inputs: in *big.Int

Ideas: try different in's around the range of the prime definding the field for NISTP-256.

##DecryptPKCS1v15SessionKey
inputs: rand io.Reader, priv *PrivateKey, ciphertext []byte, key []byte

Ideas: Reader is probably fixed, try different combinations of sizes and low/high HW for both PrivateKey and ciphertext, not sure about the key.

##decryptPKCS1v15
inputs: rand io.Reader, priv *PrivateKey, ciphertext []byte

Ideas: Almost the same as above.

##SignPKCS1v15
inputs: rand io.Reader, priv *PrivateKey, hash crypto.Hash, hashed []byte

Ideas: Similar as above, try different hashing functions.
