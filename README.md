# cloudra
Old project ( cloudra.in, when i used to own the domain ) written in go lang ( mostly initial version ) luckily still compiles.
this was inspired by transfer.sh[http://transfer.sh/], with Google AppEngine DataStore support ( then free tire; now not )
this stored is used to upload files to your bucket and it gives you short link like cloudra.in/p/<hash>

```
Usage of ./stored:
  -bucket string
    	file bucket
  -f string
    	file (default "/dev/null")
  -life string
    	life in minutes (default "200")
```
