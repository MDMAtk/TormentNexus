#!/usr/bin/env python3
"""Add Stripe checkout JS to hypernexus.site landing page"""

filepath = "/var/www/hypernexus.site/index.html"

with open(filepath, "r") as f:
    content = f.read()

stripe_js = """
<script>
function startStripeCheckout(tier) {
  var xhr = new XMLHttpRequest();
  xhr.open('POST', '/api/v1/billing/checkout', true);
  xhr.setRequestHeader('Content-Type', 'application/json');
  xhr.onload = function() {
    if (xhr.status === 200) {
      var resp = JSON.parse(xhr.responseText);
      if (resp.data && resp.data.sessionUrl) window.location.href = resp.data.sessionUrl;
      else if (resp.url) window.location.href = resp.url;
    } else {
      window.location.href = '/contact.php?tier=' + tier;
    }
  };
  xhr.onerror = function() {
    window.location.href = '/contact.php?tier=' + tier;
  };
  xhr.send(JSON.stringify({
    plan: tier,
    successUrl: window.location.origin + '/thank-you',
    cancelUrl: window.location.href
  }));
}
</script>
</body>"""

content = content.replace("</body>", stripe_js)

with open(filepath, "w") as f:
    f.write(content)

print("Stripe checkout JS added to hypernexus.site")
