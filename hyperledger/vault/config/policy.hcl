# 允許建立/讀取 指定 Transit 金鑰（建立 wrap 金鑰與簽章金鑰）
path "transit/keys/platform"            { capabilities = ["read","update"] }
path "transit/keys/platform-wrap"       { capabilities = ["read","update"] }
path "transit/keys/clinic-*-wrap"       { capabilities = ["read","update"] }
path "transit/keys/clinic-*"            { capabilities = ["read","update"] }
path "transit/keys/user-*"              { capabilities = ["read","update"] }
path "transit/keys/insurer-*"           { capabilities = ["read","update"] }

# 簽章（ECDSA）：Fabric 交易簽章
path "transit/sign/platform"            { capabilities = ["update"] }
path "transit/sign/clinic-*"            { capabilities = ["update"] }
path "transit/sign/user-*"              { capabilities = ["update"] }
path "transit/sign/insurer-*"           { capabilities = ["update"] }

# 包/解資料金鑰（AES wrap）
path "transit/encrypt/platform-wrap"    { capabilities = ["update"] }
path "transit/decrypt/platform-wrap"    { capabilities = ["update"] }
path "transit/encrypt/clinic-*-wrap"    { capabilities = ["update"] }
path "transit/decrypt/clinic-*-wrap"    { capabilities = ["update"] }

# 取得公鑰（讀取簽章金鑰公鑰）
path "transit/keys/*"                   { capabilities = ["read"] }

# KV v2 憑證存取（users/clinics/insurers/platform）
path "kv/data/users/*"                  { capabilities = ["create","update","read","delete","list"] }
path "kv/metadata/users/*"              { capabilities = ["read","list"] }
path "kv/data/clinics/*"                { capabilities = ["create","update","read","delete","list"] }
path "kv/metadata/clinics/*"            { capabilities = ["read","list"] }
path "kv/data/insurers/*"               { capabilities = ["create","update","read","delete","list"] }
path "kv/metadata/insurers/*"           { capabilities = ["read","list"] }
path "kv/data/platform"                 { capabilities = ["create","update","read","delete","list"] }
path "kv/metadata/platform"             { capabilities = ["read","list"] }