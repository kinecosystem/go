run_recipe File.dirname(__FILE__) + "/_common_accounts.rb"

use_manual_close


create_account :scott,  :master, 100
create_account :bartek, :master, 100
create_account :andrew, :master, 100

close_ledger

payment :scott, :andrew,  [:native, 5]
kp = Stellar::KeyPair.from_seed("SBQW4ZDSMV3SAIBAEAQCAIBAEAQCAIBAEAQCAIBAEAQCAIBAEAQCA65I") # andrew's account
add_signer :scott, kp, 1

close_ledger

