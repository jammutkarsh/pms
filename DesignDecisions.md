# Why this not that

## Why Decode-Encode V/S Marshal-Unmarshal

Source: <https://blog.devgenius.io/to-unmarshal-or-to-decode-json-processing-in-go-explained-e92fab5b648f>

According to this blog, the encode and decode methods deal with streams of data.
That is, the entire data is not loaded into memory at once.
While Marshal and Unmarshal methods load the entire data into the memory,

From a performance standpoint, they both are equivalent. But from a learning standpoint, I know how to use the marshal and unmarshal methods. So this could be an opportunity to learn the encode and decode methods.

## Challange: Debugging TUI

A TUI always involves `STDIN` & `STDOUT`, hence it's not possible to debug in the typical GUI editor like VScode. Since the debug console is not interactive.

**Solution**: <https://www.jjshanks.net/posts/go-stdin-dlv/>

```bash
# In Terminal 1
dlv debug --headless --listen :4747 .

# In Terminal 2
dlv connect :4747
```

Now you can start using the TUI in **Terminal 1** and debug each step in Terminal 2 using `dlv debug` commands.
