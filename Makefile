DRAFT = draft-nandakumar-agentproto-trail.md
LIBDIR := ../moq-transport/lib
include $(LIBDIR)/main.mk

# Preserve XML files - do not delete them as intermediates or during clean
.PRECIOUS: %.xml
.SECONDARY: %.xml
