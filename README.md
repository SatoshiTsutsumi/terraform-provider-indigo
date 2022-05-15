# Terraform Provider Indigo (alpha)
Terraform Provider for WebARENA Indigo

# Prerequisites
- [indigo-client-go](https://github.com/UndefxDev/indigo-client-go)
- [terraform-plugin-sdk v2](https://github.com/hashicorp/terraform-plugin-sdk)

# Support Resources
- SSH Key
- Instance
- Snapshot
- Firewall

# How to use
```
$ make
```

Please refer to example/main.tf.
Values for each parameter are under investigation as they seem NOT documented.

## Note
WebARENA Indigo API offers very limited functionality and unfriendly UX to deploy commercial systems.
Terraform Provider Indigo is just an experimental trial unless the API provides commercial grade functionaly and performance.

Known limitations are listed below:
- Rate limit (< 1 call/s)
- Call limit for creating instances (< around 10 instances/day)

## Workaround
If you encounter an error such as "Too Many Request" while doing "terraform apply" , please retry it.

## References
- [WebARENA Indigo API](https://indigo.arena.ne.jp/userapi/)
