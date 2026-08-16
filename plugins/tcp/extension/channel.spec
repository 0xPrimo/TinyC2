x64:
	push $OBJECT
		make pic +optimize +gofirst

		run "services.spec"

		mergelib "./libs/libtcg/libtcg.x64.zip"

		push $CONFIG
		    preplen
		    link "config"

		export