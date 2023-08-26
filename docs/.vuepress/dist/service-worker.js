/**
 * Welcome to your Workbox-powered service worker!
 *
 * You'll need to register this file in your web app and you should
 * disable HTTP caching for this file too.
 * See https://goo.gl/nhQhGp
 *
 * The rest of the code is auto-generated. Please don't update this file
 * directly; instead, make changes to your Workbox build configuration
 * and re-run your build process.
 * See https://goo.gl/2aRDsh
 */

importScripts("https://storage.googleapis.com/workbox-cdn/releases/4.3.1/workbox-sw.js");

self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting();
  }
});

/**
 * The workboxSW.precacheAndRoute() method efficiently caches and responds to
 * requests for URLs in the manifest.
 * See https://goo.gl/S9QRab
 */
self.__precacheManifest = [
  {
    "url": "404.html",
    "revision": "5c85637f597332be951527b070edf10f"
  },
  {
    "url": "assets/css/0.styles.fbd442e9.css",
    "revision": "7b6607f23dcfe043c879f4340ae4b25c"
  },
  {
    "url": "assets/img/search.83621669.svg",
    "revision": "83621669651b9a3d4bf64d1a670ad856"
  },
  {
    "url": "assets/js/1.3e70dec9.js",
    "revision": "4733ce017504b3fc7aefe462fc5d046d"
  },
  {
    "url": "assets/js/12.4cd07bc5.js",
    "revision": "0608c4652ad8f05fa8d44dd5ca578f51"
  },
  {
    "url": "assets/js/13.b4d09a7e.js",
    "revision": "8a0e23e7a7681ca33544b1850e3d26e9"
  },
  {
    "url": "assets/js/14.4b359ba2.js",
    "revision": "de3fe578fa5485f53c9bdae2e669a4dd"
  },
  {
    "url": "assets/js/15.116c4bfc.js",
    "revision": "7bb8fcd1d396cda82ac9de3b7c45b37d"
  },
  {
    "url": "assets/js/16.7b87341d.js",
    "revision": "e862fe7cfe64b2886c2567fc0e1b11c3"
  },
  {
    "url": "assets/js/17.2cde4da0.js",
    "revision": "12366022fc36c59c6004bfe5c34f59f8"
  },
  {
    "url": "assets/js/18.53ff10d3.js",
    "revision": "b45f7a46bbc4f8f482575604a429e9be"
  },
  {
    "url": "assets/js/19.5831e1ee.js",
    "revision": "9d9a2da5b94b10df7721b5e205cc80c6"
  },
  {
    "url": "assets/js/2.62b9495e.js",
    "revision": "a6bf8cf231a8b1f810c363602384f1c1"
  },
  {
    "url": "assets/js/20.4f8319b4.js",
    "revision": "8d8407fb723fe0cd23be51d393d53381"
  },
  {
    "url": "assets/js/21.56697f36.js",
    "revision": "3bcd1c641f816147e1ae7cf26dc5e444"
  },
  {
    "url": "assets/js/22.6f2c8690.js",
    "revision": "dd5423ff071b5a9c5736d53f5753acf5"
  },
  {
    "url": "assets/js/23.7d243eaf.js",
    "revision": "b3886cafa0d9964c4fcf3d7419a28473"
  },
  {
    "url": "assets/js/24.7272cb65.js",
    "revision": "120b3b05ed5258ef19d49ba89cf030be"
  },
  {
    "url": "assets/js/25.83814250.js",
    "revision": "33910696000ece73c28634edda9b48e3"
  },
  {
    "url": "assets/js/26.3a5c06cb.js",
    "revision": "ab67841682a9593e4fb103afecd05072"
  },
  {
    "url": "assets/js/27.57bfc038.js",
    "revision": "7034d82994d9a3094d0b1f8aa73009f7"
  },
  {
    "url": "assets/js/28.0fa0a592.js",
    "revision": "9ccdb7c67ddda9239e451c4dd8702573"
  },
  {
    "url": "assets/js/29.7ae92863.js",
    "revision": "5aa76b4c22047633544896929827c7ca"
  },
  {
    "url": "assets/js/3.a84da161.js",
    "revision": "4bcee03349f01486cde93cf6de623985"
  },
  {
    "url": "assets/js/30.171f9842.js",
    "revision": "540f3343616b9dd65cf7c87cd6b9271f"
  },
  {
    "url": "assets/js/31.665ebcae.js",
    "revision": "85beb07c9f206e17e55dcdc0fec89c83"
  },
  {
    "url": "assets/js/32.e2f69e92.js",
    "revision": "78c8177447368fbb335822c4dd4d332b"
  },
  {
    "url": "assets/js/4.7668b4a8.js",
    "revision": "2217964c016de0911a04e982c9dcb835"
  },
  {
    "url": "assets/js/5.2dd90f13.js",
    "revision": "b6c441b9a76043bf5538243cb6f9b8f8"
  },
  {
    "url": "assets/js/6.99df0be8.js",
    "revision": "fbdf671951604cd00685d9f2ecf18c30"
  },
  {
    "url": "assets/js/7.82964bba.js",
    "revision": "07eb7e88994ad6e16532cfdbb7972dce"
  },
  {
    "url": "assets/js/8.e2ee25a8.js",
    "revision": "8ff4ccdbe82cafa9e0561bc810baa2e8"
  },
  {
    "url": "assets/js/app.0df30dd7.js",
    "revision": "451d9a4413643dbddf53ece8497efdda"
  },
  {
    "url": "assets/js/vendors~docsearch.19b85b09.js",
    "revision": "36558bee2bc75cc2044cc1424f4f14c0"
  },
  {
    "url": "assets/js/vendors~flowchart.1c686b1f.js",
    "revision": "98d2b632e5d9a4f7f8abfb6a8fe999c8"
  },
  {
    "url": "cmd/dsg.html",
    "revision": "a7b9c283bced53ddadacd80dc910a3a9"
  },
  {
    "url": "cmd/gmm.html",
    "revision": "924d6099cd083fa701a2c3e10fcf01ca"
  },
  {
    "url": "cmd/index.html",
    "revision": "673da08c08f3458acecc3072b7d9a34d"
  },
  {
    "url": "cmd/jpu.html",
    "revision": "ae185fc50684c14cdc0cd91e8daadc99"
  },
  {
    "url": "index.html",
    "revision": "c46b7cedc777d99250aecbee3f3fd2ea"
  },
  {
    "url": "utils/index.html",
    "revision": "49f77aa247affe58e1e5b457b7044f2f"
  }
].concat(self.__precacheManifest || []);
workbox.precaching.precacheAndRoute(self.__precacheManifest, {});
addEventListener('message', event => {
  const replyPort = event.ports[0]
  const message = event.data
  if (replyPort && message && message.type === 'skip-waiting') {
    event.waitUntil(
      self.skipWaiting().then(
        () => replyPort.postMessage({ error: null }),
        error => replyPort.postMessage({ error })
      )
    )
  }
})
