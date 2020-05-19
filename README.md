## keep your secret

### why did I need this?

Do you have secret, which you would share for nobody? I have.

For nobody means, you wouldn't share it with your father, mother, your boy friend, girl friend, even yourself.

Yes, even yourself. Maybe we have some bullshit, and write it down just for forgetting it. We won't check it anymore.

Then here it is, if you have the same require, maybe dong would help you.

### which level of secret?

Dong has no account system, it just has an `abracadabra`, it's the path to your secret.

Every secret could contain some papers. Paper is a piece of secret.

There is a so big abracadabra space, which means you have a lot of abracadabra for using. Some for password, and some for story.

Brower will encrypt your words and then puts it to the server, the server also encrypt your data.

There is no sequence order field, which means no one could tell the abracadabra count, paper count by database.

### how to achieve this?

On browser client, we use sha256 and crypto to encrypt your abracadabra and your data.

At the backend, we use the encrypted abracadabra and salt to encrypt data id and your data.

The algo is here:

```
              sha256            split
abracadabra ----------> s256 ------------> l128 | r128

       aes (l128)
data -------------> encryptedData

       sha256 (order)
r128 ------------------> id

       sha256 ( + salt)
r128 --------------------> aesKey

                aes (aesKey)
encryptedData ----------------> encryptedEncryptedData
```

### install

Install dong is simple, clone this repo and run the following command:

```shell script
# build from source
go build

# start server
./dong -addr :5300

# then open http://localhost:5300/app/ from your browser
```