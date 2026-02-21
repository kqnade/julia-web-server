(function () {
  "use strict";

  var TILE_SIZE = 256;
  var canvas = document.getElementById("canvas");
  var ctx = canvas.getContext("2d");
  var errorEl = document.getElementById("error");
  var btn = document.getElementById("generate");

  btn.addEventListener("click", generate);

  function val(id) {
    return document.getElementById(id).value.trim();
  }

  function generate() {
    errorEl.textContent = "";
    btn.disabled = true;

    var minX = parseFloat(val("min_x"));
    var maxX = parseFloat(val("max_x"));
    var minY = parseFloat(val("min_y"));
    var maxY = parseFloat(val("max_y"));
    var cReal = parseFloat(val("c_real"));
    var cImag = parseFloat(val("c_imag"));

    if (!Number.isFinite(minX) || !Number.isFinite(maxX) ||
        !Number.isFinite(minY) || !Number.isFinite(maxY) ||
        !Number.isFinite(cReal) || !Number.isFinite(cImag)) {
      errorEl.textContent = "All parameters must be valid finite numbers.";
      btn.disabled = false;
      return;
    }

    var canvasW = canvas.width;
    var canvasH = canvas.height;
    var cols = Math.ceil(canvasW / TILE_SIZE);
    var rows = Math.ceil(canvasH / TILE_SIZE);

    // Clear canvas
    ctx.fillStyle = "#000";
    ctx.fillRect(0, 0, canvasW, canvasH);

    var fetches = [];

    for (var row = 0; row < rows; row++) {
      for (var col = 0; col < cols; col++) {
        (function (col, row) {
          var pxLeft = col * TILE_SIZE;
          var pxTop = row * TILE_SIZE;

          // Tile dimensions (clamp at canvas edge)
          var tileW = Math.min(TILE_SIZE, canvasW - pxLeft);
          var tileH = Math.min(TILE_SIZE, canvasH - pxTop);

          // Map tile pixel range to complex plane
          var tMinX = minX + (maxX - minX) * pxLeft / canvasW;
          var tMaxX = minX + (maxX - minX) * (pxLeft + tileW) / canvasW;
          var tMinY = minY + (maxY - minY) * pxTop / canvasH;
          var tMaxY = minY + (maxY - minY) * (pxTop + tileH) / canvasH;

          var url =
            "/satori/julia/api" +
            "?min_x=" + tMinX +
            "&max_x=" + tMaxX +
            "&min_y=" + tMinY +
            "&max_y=" + tMaxY +
            "&comp_const=" + encodeURIComponent(cReal + "," + cImag) +
            "&width=" + tileW +
            "&height=" + tileH;

          var p = fetch(url).then(function (resp) {
            var ct = resp.headers.get("Content-Type") || "";
            if (!resp.ok) {
              if (ct.indexOf("application/json") !== -1) {
                return resp.json().then(function (data) {
                  throw new Error(data.error || "Unknown server error");
                });
              }
              throw new Error("Server error: " + resp.status + " " + resp.statusText);
            }
            return resp.arrayBuffer();
          }).then(function (ab) {
            var floats = new Float32Array(ab);
            var imgData = ctx.createImageData(tileW, tileH);
            var pixels = imgData.data;

            for (var i = 0; i < floats.length; i++) {
              var smooth = floats[i];
              var off = i * 4;

              if (smooth < 0) {
                // Interior: black
                pixels[off] = 0;
                pixels[off + 1] = 0;
                pixels[off + 2] = 0;
              } else {
                // HSV coloring
                var hue = (smooth * 10) % 360;
                var rgb = hsvToRgb(hue, 1.0, 1.0);
                pixels[off] = rgb[0];
                pixels[off + 1] = rgb[1];
                pixels[off + 2] = rgb[2];
              }
              pixels[off + 3] = 255; // alpha
            }

            ctx.putImageData(imgData, pxLeft, pxTop);
          });

          fetches.push(p);
        })(col, row);
      }
    }

    Promise.allSettled(fetches).then(function (results) {
      var errors = [];
      results.forEach(function (result) {
        if (result.status === "rejected") {
          errors.push(result.reason.message);
        }
      });
      if (errors.length > 0) {
        errorEl.textContent = errors[0];
      }
      btn.disabled = false;
    });
  }

  // Convert HSV to RGB. h in [0,360), s and v in [0,1].
  // Returns [r, g, b] each in [0,255].
  function hsvToRgb(h, s, v) {
    var c = v * s;
    var x = c * (1 - Math.abs(((h / 60) % 2) - 1));
    var m = v - c;
    var r, g, b;

    if (h < 60) {
      r = c; g = x; b = 0;
    } else if (h < 120) {
      r = x; g = c; b = 0;
    } else if (h < 180) {
      r = 0; g = c; b = x;
    } else if (h < 240) {
      r = 0; g = x; b = c;
    } else if (h < 300) {
      r = x; g = 0; b = c;
    } else {
      r = c; g = 0; b = x;
    }

    return [
      Math.round((r + m) * 255),
      Math.round((g + m) * 255),
      Math.round((b + m) * 255)
    ];
  }
})();
