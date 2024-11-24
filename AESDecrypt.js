
async function AESDecrypt(cipher, password) {
    let parts = cipher.split("|");
    let ciphertext = parts[1];
    let nonce = parts[0];
    const ciphertextBuffer = hexToBytes(ciphertext)
    let encoder = new TextEncoder();
    let data = encoder.encode(password);

    let hashBuffer = await crypto.subtle.digest('SHA-256', data);

    let hash = Array.from(new Uint8Array(hashBuffer));

    const hashKey = new Uint8Array(hash);
    const key = await window.crypto.subtle.importKey(
        'raw',
        hashKey, {
        name: 'AES-GCM',
    },
        false,
        ['decrypt']
    )
    let iv = hexToBytes(nonce)
    const decrypted = await window.crypto.subtle.decrypt({
        name: 'AES-GCM',
        iv: iv,
        tagLength: 128,
    },
        key,
        new Uint8Array(ciphertextBuffer)
    )
    return new TextDecoder('utf-8').decode(new Uint8Array(decrypted))
}
function hexToBytes(hexString) {
    if (hexString.startsWith("0x") || hexString.startsWith("0X")) {
        hexString = hexString.slice(2);
    }

    const bytes = new Uint8Array(hexString.length / 2);
    for (let i = 0; i < hexString.length; i += 2) {
        bytes[i / 2] = parseInt(hexString.substr(i, 2), 16);
    }

    return bytes;
}
console.log("js load");
let title = document.title
if (localStorage.getItem(title)!== null) {
    decryption(localStorage.getItem(title))
}
const submitButton = document.getElementById('secret-submit');
submitButton.addEventListener('click', function (event) {
    event.preventDefault(); // Blocking the default form submission behavior
    checkPassword();
});

function checkPassword() {
    const passwordInput = document.querySelector('input[name="password"]');
    const password = passwordInput.value;
    decryption(password)
}
function decryption(password) {
    let secretElement = document.getElementById('secret');
    let ciphertext = secretElement.innerText;
    AESDecrypt(ciphertext, password).then(plaintext => {
        document.getElementById("verification").style.display = "none";
        let verificationElement = document.getElementById('verification');
        let htmlText =  marked.parse(plaintext);
        verificationElement.insertAdjacentHTML('afterend', htmlText);
        if (localStorage.getItem(title) !==password)localStorage.setItem(title, password);
    }).catch(error => {
        let repOnFail = "Incorrect password!!"
        let repOnFailElement = document.getElementById('repOnFail');
        if (repOnFailElement != null && repOnFailElement.hasAttribute("text")) {
            repOnFail = repOnFailElement.getAttribute("text")
        }
        alert(repOnFail);
        console.error("Failed to decrypt",error);
    });
}