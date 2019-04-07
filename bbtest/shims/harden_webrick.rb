require 'openssl'
require 'webrick'
require 'webrick/ssl'

module WEBrick

  class GenericServer
    def setup_ssl_context(config)
      unless config[:SSLCertificate]

        key = OpenSSL::PKey::RSA.new(1024)
        subject = "/C=CZ/ST=Czechia/L=Prague/O=OpenBanking/OU=IT/CN=localhost/emailAddress=jan.cajthaml@gmail.com"

        cert = OpenSSL::X509::Certificate.new
        cert.subject = cert.issuer = OpenSSL::X509::Name.parse(subject)
        cert.not_before = Time.now
        cert.not_after = Time.now + 30 * 60
        cert.public_key = key.public_key
        cert.serial = 0x0
        cert.version = 0x1

        ef = OpenSSL::X509::ExtensionFactory.new()
        ef.subject_certificate = cert
        ef.issuer_certificate = cert
        cert.extensions = [
          ef.create_extension("basicConstraints","CA:TRUE", true),
          ef.create_extension("subjectKeyIdentifier", "hash"),
          ef.create_extension("keyUsage", "cRLSign,digitalSignature,keyCertSign", true),
          ef.create_extension("extendedKeyUsage", "serverAuth", true)
        ]

        cert.add_extension ef.create_extension("authorityKeyIdentifier", "keyid:always,issuer:always")

        cert.sign key, OpenSSL::Digest::SHA256.new()

        config[:SSLCertificate] = cert
        config[:SSLPrivateKey] = key
        config[:SSLCiphers] = 'ALL:!aNULL:!eNULL:!SSLv2'
        config[:SSLVerifyClient] = OpenSSL::SSL::VERIFY_PEER
        config[:SSLOptions] = OpenSSL::SSL::OP_NO_SSLv2 + OpenSSL::SSL::OP_NO_SSLv3
        config[:SSLVersion] = :TLSv1_2
        config[:SSLStartImmediately] = true
      end

      ctx = OpenSSL::SSL::SSLContext.new()
      ctx.set_params({
        key: config[:SSLPrivateKey],
        cert: config[:SSLCertificate],
        verify_mode: config[:SSLVerifyClient],
        timeout: config[:SSLTimeout],
        options: config[:SSLOptions],
        ciphers: config[:SSLCiphers],
        ssl_version: config[:SSLVersion]
      })
      ctx.verify_mode = config[:SSLVerifyClient]
      ctx.ssl_version = config[:SSLVersion]
      ctx
    end
  end

end
