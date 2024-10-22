import os
from shutil import which
from secrets import token_bytes
from base64 import b64encode
from cryptography import x509
from cryptography.x509.oid import NameOID
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import rsa
from datetime import datetime, timedelta


ROOT_DIR = os.path.dirname(__file__)
KITSUNEC2_DIR = "kitsuneC2"
KITSUNEC2_SERVER_DIR = "kitsuneC2/server"
GO_SERVER_ENV = "kitsuneC2/server/.template-env"
KITSUNEC2_FRONTEND_ENV = "kitsune-frontend/.template-env"
FRONTEND_CERTIFICATE_DIR = "kitsune-frontend/certificates"
BACKEND_CERTIFICATE_DIR = "kitsuneC2/certificates"
CERTIFICATE_NAME = "kc2SSL"

def main():
    if os.name != "posix":
        print("Unsupported operating system, quitting! For windows, CLI install is available by checking the \"releases\" tab on Github")
        exit(1)

    if not webDependenciesInstalled():
        print("Docker is not installed, or not in $PATH. Please install docker and add it to $PATH")
        exit(1)
    linuxWebInstall()


def linuxWebInstall():
    username = getInput("Username for web-login: ")
    password = getInput("Password for web-login: ")

    bindAdress = getInput("Domain (or ip address) that server will be running on (e.g. \"kitsunec2.com\" or \"123.456.789.123\" or \"localhost\"): ")
    
    if os.path.isfile(os.path.join(FRONTEND_CERTIFICATE_DIR, CERTIFICATE_NAME + ".key")) and os.path.isfile(os.path.join(FRONTEND_CERTIFICATE_DIR, CERTIFICATE_NAME + ".pem")):
        print("Existing frontend certificate files found, skipping SSL-certificate generation...\n")
    else:
        print(f"Generating frontend SSL certificate (CN={bindAdress}, valid for {999} days)...\n")
        generateSSLCert(FRONTEND_CERTIFICATE_DIR, bindAdress, 999)
        print(f"Written frontend SSL certificate to: {os.path.abspath(FRONTEND_CERTIFICATE_DIR)}\n")

    if os.path.isfile(os.path.join(BACKEND_CERTIFICATE_DIR, CERTIFICATE_NAME + ".key")) and os.path.isfile(os.path.join(BACKEND_CERTIFICATE_DIR, CERTIFICATE_NAME + ".pem")):
        print("Existing backend certificate files found, skipping SSL-certificate generation...\n")
    else:
        print(f"Generating backend SSL certificate (CN={bindAdress}, valid for {999} days)...\n")
        generateSSLCert(BACKEND_CERTIFICATE_DIR, bindAdress, 999)
        print(f"Written backend SSL certificate to: {os.path.abspath(BACKEND_CERTIFICATE_DIR)}\n")


    print("Generating authentication secrets...\n")
    nextAuthSecret = b64encode(token_bytes(32)).decode()
    apiAuthToken = b64encode(token_bytes(32)).decode()

    print("Writing frontend .env file...\n")
    redirectUrl = "https://" + bindAdress + ":8080/"
    envVars = ""
    with open(os.path.join(ROOT_DIR, KITSUNEC2_FRONTEND_ENV), 'r') as file:
        envVars = file.read()
        envVars = envVars.format(_apiAuthToken=apiAuthToken,_nextAuthUrl=redirectUrl, _nextAuthSecret=nextAuthSecret, _username=username, _password=password)
        file.close()
    with open(os.path.join(ROOT_DIR, KITSUNEC2_FRONTEND_ENV.replace(".template-env", ".env")), 'w') as file:
        file.write(envVars)
        file.close()
    
    print("Writing backend .env file...\n")
    with open(os.path.join(ROOT_DIR, GO_SERVER_ENV), 'r') as file:
        envVars = file.read()
        envVars = envVars.format(_apiAuthToken=apiAuthToken,_webEnabled="true")
        file.close()
    with open(os.path.join(ROOT_DIR, GO_SERVER_ENV.replace(".template-env", ".env")), 'w') as file:
        file.write(envVars)
        file.close()


    print("Environment setup done. Navigate to {} and run command \"sudo docker compose up\"".format(ROOT_DIR))
    

def generateSSLCert(outDir, cn, validDays):
    # Set up parameters for certificate and key
    common_name = cn
    valid_days = validDays

    # Generate RSA private key
    private_key = rsa.generate_private_key(
        public_exponent=65537,
        key_size=2048,
        backend=default_backend()
    )

    # Create a certificate signing request (CSR)
    subject = issuer = x509.Name([
        x509.NameAttribute(NameOID.COMMON_NAME, common_name),
    ])

    # You may also add a SAN for compatibility with modern browsers
    certificate = (
        x509.CertificateBuilder()
        .subject_name(subject)
        .issuer_name(issuer)
        .public_key(private_key.public_key())
        .serial_number(x509.random_serial_number())
        .not_valid_before(datetime.now())
        .not_valid_after(datetime.now() + timedelta(days=valid_days))
        .add_extension(x509.SubjectAlternativeName([x509.DNSName(common_name)]), critical=False)
        .sign(private_key, hashes.SHA256(), default_backend())
    )

    # Write the private key to a file
    with open(os.path.join(outDir, "kc2SSL.key"), "wb") as f:
        f.write(private_key.private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.TraditionalOpenSSL,
            encryption_algorithm=serialization.NoEncryption(),
        ))

    
    # Write the certificate to a file
    with open(os.path.join(outDir, "kc2SSL.pem"), "wb") as f:
        f.write(certificate.public_bytes(serialization.Encoding.PEM))

def webDependenciesInstalled():
    return which("docker") is not None

def getInput(prompt):
    i = input(prompt + "\n")
    return i.lower()

if __name__ == "__main__":
    main()










