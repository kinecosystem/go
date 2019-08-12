run_recipe File.dirname(__FILE__) + "/_common_accounts.rb"

use_manual_close


create_account :scott,  :master, 100
create_account :bartek, :master, 100
create_account :andrew, :master, 100

# create accounts:
# kp1
# GCJ6M7DRW5RZW73UEHEOGIGKLY5BWA5QNR6WX25G5BSYWDQYR77DSIH7
# SAYYMGTGIWMR6VO243LU5GG2J5YYMAPFK5ZSTITHPUV5DVUS3K6Q4ZNF

# kp2
# GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y
# SBMTHT5NRNVJFBYDON5COPZRCRCMDOQNRPAOMVJJV6W7DFW32EKRN6RV

# kp3
# GBSD7E6QRGGEBFWVF4VNSAWPQJW7FFLYEV7O3DZUIRLYATFMLCRJGSNQ
# SBOTI576E6NH4DWMRSWFW2SWMY7CCCGUX33DDJCG4RZLDSJ6RKZHI7NB

# kp4
# GBS2EWY4PAVG4SOSGZ6CM5F4SIJUAIEXGKXEXCGTMCAK6FFVDDDXVHI6
# SC3CILTVMLKPF7YJXXIESDXNIA2QX6XX47WMLAY3UI5YQL75Y2SRXOZN

# kp5
# GBKNDEBYVKMVRRR376KGV77XMEYKOHKZQRN5TEOKJYJZI3VHBV7YKLJZ
# SDB7GWKMZ4Q2DLXNPW3KGDQZEKIXVGFQR67VF73ZDZZ7SWMUX5ADIHVN

# kp6
# GCBJRUEBARNP5HLASDN3ZNVQSYFFO2ETXCLZVCI2DSHTFBMA6R37U6S2
# SCKKSSNCJ4MP2373BQM44LPDFYET2LTDJOBI4P66264TQ6PO4UBCUZR7

# kp7
# GDSLCGMN4WK2SAOANYYOSIATOT2CWTSM5AHBQYAAXXJYEHI5CDEHRYIL
# SAQOQ2P772O6MJIVOSDSLNNOWQDTP6RZXWX4AXY6VTGU5GXUJTXZTSC4

# kp8
# GCDI34A4Y4ELTHGILNEC32NOAEIJLFRJFMJ75PFJAV7VXFU6BPQVYTEN
# SARTM2RF4YXSTU672RXHF2HPAX7Y3WLEDZV52S3ZLDVBM6KMTYLOIUU6

# kp9
# GCXOZEDULG6VUL7RMUPZDAMFAZ5WGXZDYNOGYLFK3YAPH2PIYGHSJXPT
# SAQW5D2IB4N7JL2LIO6UGEB3TQBTWUSRREGADFJBKQMKJ3UH3OWXAD2E

account :kp1, Stellar::KeyPair.from_seed("SAYYMGTGIWMR6VO243LU5GG2J5YYMAPFK5ZSTITHPUV5DVUS3K6Q4ZNF")
account :kp2, Stellar::KeyPair.from_seed("SBMTHT5NRNVJFBYDON5COPZRCRCMDOQNRPAOMVJJV6W7DFW32EKRN6RV")
account :kp3, Stellar::KeyPair.from_seed("SBOTI576E6NH4DWMRSWFW2SWMY7CCCGUX33DDJCG4RZLDSJ6RKZHI7NB")
account :kp4, Stellar::KeyPair.from_seed("SC3CILTVMLKPF7YJXXIESDXNIA2QX6XX47WMLAY3UI5YQL75Y2SRXOZN")

account :kp5, Stellar::KeyPair.from_seed("SDB7GWKMZ4Q2DLXNPW3KGDQZEKIXVGFQR67VF73ZDZZ7SWMUX5ADIHVN")
account :kp6, Stellar::KeyPair.from_seed("SCKKSSNCJ4MP2373BQM44LPDFYET2LTDJOBI4P66264TQ6PO4UBCUZR7")
account :kp7, Stellar::KeyPair.from_seed("SAQOQ2P772O6MJIVOSDSLNNOWQDTP6RZXWX4AXY6VTGU5GXUJTXZTSC4")

account :kp8, Stellar::KeyPair.from_seed("SARTM2RF4YXSTU672RXHF2HPAX7Y3WLEDZV52S3ZLDVBM6KMTYLOIUU6")
account :kp9, Stellar::KeyPair.from_seed("SAQW5D2IB4N7JL2LIO6UGEB3TQBTWUSRREGADFJBKQMKJ3UH3OWXAD2E")


# accounts must be created in a ledger prior to actual manipulation.
create_account :kp1, :master, 100
create_account :kp2, :master, 200
create_account :kp3, :master, 300
create_account :kp4, :master, 400

create_account :kp5, :master, 500
create_account :kp6, :master, 600
create_account :kp7, :master, 700

create_account :kp8, :master, 800
create_account :kp9, :master, 900

close_ledger

payment :scott, :andrew,  [:native, 5]

# add andrew as signer on sopott
kp = Stellar::KeyPair.from_seed("SBQW4ZDSMV3SAIBAEAQCAIBAEAQCAIBAEAQCAIBAEAQCAIBAEAQCA65I") # andrew's account
add_signer :scott, kp, 1

# add kp2 as signer in kp1, kp1 pays for this operation
kp2 = Stellar::KeyPair.from_seed("SBMTHT5NRNVJFBYDON5COPZRCRCMDOQNRPAOMVJJV6W7DFW32EKRN6RV")
add_signer :kp1, kp2, 1

# add kp3 as signer in kp2
kp3 = Stellar::KeyPair.from_seed("SBOTI576E6NH4DWMRSWFW2SWMY7CCCGUX33DDJCG4RZLDSJ6RKZHI7NB")
add_signer :kp2, kp3, 1

# add kp4 as signer in kp3
kp4 = Stellar::KeyPair.from_seed("SC3CILTVMLKPF7YJXXIESDXNIA2QX6XX47WMLAY3UI5YQL75Y2SRXOZN") # kp4
add_signer :kp3, kp4, 1


close_ledger

# add kp5 as signer in kp6
kp5 = Stellar::KeyPair.from_seed("SDB7GWKMZ4Q2DLXNPW3KGDQZEKIXVGFQR67VF73ZDZZ7SWMUX5ADIHVN") # kp5
add_signer :kp6, kp5, 1

# add kp6 as signer in kp7
kp6 = Stellar::KeyPair.from_seed("SCKKSSNCJ4MP2373BQM44LPDFYET2LTDJOBI4P66264TQ6PO4UBCUZR7") # kp6
add_signer :kp7, kp6, 1

# add kp7 as signer in kp5
kp7 = Stellar::KeyPair.from_seed("SAQOQ2P772O6MJIVOSDSLNNOWQDTP6RZXWX4AXY6VTGU5GXUJTXZTSC4") # kp7
add_signer :kp5, kp7, 1

close_ledger

# kp8->kp9->kp8
# add kp as signer in kp5
kp8 = Stellar::KeyPair.from_seed("SARTM2RF4YXSTU672RXHF2HPAX7Y3WLEDZV52S3ZLDVBM6KMTYLOIUU6") # kp8
add_signer :kp9, kp8, 1

kp9 = Stellar::KeyPair.from_seed("SAQW5D2IB4N7JL2LIO6UGEB3TQBTWUSRREGADFJBKQMKJ3UH3OWXAD2E") # kp9
add_signer :kp8, kp9, 1

close_ledger
