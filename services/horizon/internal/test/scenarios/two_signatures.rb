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


account :kp1, Stellar::KeyPair.from_seed("SAYYMGTGIWMR6VO243LU5GG2J5YYMAPFK5ZSTITHPUV5DVUS3K6Q4ZNF")
account :kp2, Stellar::KeyPair.from_seed("SBMTHT5NRNVJFBYDON5COPZRCRCMDOQNRPAOMVJJV6W7DFW32EKRN6RV")
account :kp3, Stellar::KeyPair.from_seed("SBOTI576E6NH4DWMRSWFW2SWMY7CCCGUX33DDJCG4RZLDSJ6RKZHI7NB")
account :kp4, Stellar::KeyPair.from_seed("SC3CILTVMLKPF7YJXXIESDXNIA2QX6XX47WMLAY3UI5YQL75Y2SRXOZN")

# accounts must be created in a ledger prior to actual manipulation.
create_account :kp1, :master, 100
create_account :kp2, :master, 200
create_account :kp3, :master, 300
create_account :kp4, :master, 400

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

